package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"dailytrackr/ai-service/models"
	"dailytrackr/shared/config"
)

type GeminiService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// GeminiRequest represents request to Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents response from Gemini API
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// NewGeminiService creates a new Gemini service instance
func NewGeminiService(cfg *config.Config) *GeminiService {
	return &GeminiService{
		apiKey:  cfg.GeminiAPIKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateDailySummary generates AI-powered daily summary
func (g *GeminiService) GenerateDailySummary(activities []models.Activity) (string, error) {
	// Build activity summary for prompt
	var activityDetails strings.Builder
	totalTime := 0
	totalCost := 0

	activityDetails.WriteString("Aktivitas hari ini:\n")
	for i, activity := range activities {
		activityDetails.WriteString(fmt.Sprintf("%d. %s (%d menit)",
			i+1, activity.Title, activity.DurationMins))

		if activity.Cost != nil && *activity.Cost > 0 {
			activityDetails.WriteString(fmt.Sprintf(" - Biaya: Rp%d", *activity.Cost))
			totalCost += *activity.Cost
		}

		if activity.Note != "" {
			activityDetails.WriteString(fmt.Sprintf(" - Catatan: %s", activity.Note))
		}

		activityDetails.WriteString("\n")
		totalTime += activity.DurationMins
	}

	activityDetails.WriteString(fmt.Sprintf("\nTotal waktu: %d menit (%.1f jam)",
		totalTime, float64(totalTime)/60.0))
	activityDetails.WriteString(fmt.Sprintf("\nTotal pengeluaran: Rp%d", totalCost))

	prompt := fmt.Sprintf(`
Buatkan ringkasan harian yang menarik dan motivational berdasarkan aktivitas berikut:

%s

Tugas:
1. Buat ringkasan singkat (2-3 kalimat) tentang pencapaian hari ini
2. Berikan insights tentang produktivitas 
3. Tambahkan motivasi untuk hari berikutnya
4. Gunakan bahasa Indonesia yang friendly dan encouraging
5. Fokus pada aspek positif dan progress yang telah dibuat

Format ringkasan dalam paragraf yang mengalir, bukan poin-poin.
`, activityDetails.String())

	return g.callGeminiAPI(prompt)
}

// GenerateHabitRecommendation generates AI-powered habit recommendations
func (g *GeminiService) GenerateHabitRecommendation(activities []models.Activity, existingHabits []models.Habit) (string, error) {
	// Analyze activity patterns
	activityTypes := make(map[string]int)
	timePatterns := make(map[int]int) // hour -> count
	totalDuration := 0

	for _, activity := range activities {
		// Categorize activities (simplified)
		category := categorizeActivity(activity.Title)
		activityTypes[category]++

		hour := activity.StartTime.Hour()
		timePatterns[hour]++
		totalDuration += activity.DurationMins
	}

	// Build existing habits summary
	var existingHabitsStr strings.Builder
	existingHabitsStr.WriteString("Habit yang sudah ada:\n")
	for _, habit := range existingHabits {
		existingHabitsStr.WriteString(fmt.Sprintf("- %s (%s, progress: %d%%)\n",
			habit.Title, habit.Status, habit.Progress))
	}

	// Build activity analysis
	var activityAnalysis strings.Builder
	activityAnalysis.WriteString("Analisis aktivitas 7 hari terakhir:\n")
	activityAnalysis.WriteString(fmt.Sprintf("- Total aktivitas: %d\n", len(activities)))
	activityAnalysis.WriteString(fmt.Sprintf("- Total waktu: %.1f jam\n", float64(totalDuration)/60.0))

	// Most common activity type
	maxCount := 0
	mostCommonType := ""
	for actType, count := range activityTypes {
		if count > maxCount {
			maxCount = count
			mostCommonType = actType
		}
	}
	activityAnalysis.WriteString(fmt.Sprintf("- Jenis aktivitas terbanyak: %s (%d kali)\n",
		mostCommonType, maxCount))

	prompt := fmt.Sprintf(`
Sebagai AI productivity coach, analisis pola aktivitas user dan berikan rekomendasi habit baru yang personal dan actionable.

%s

%s

Tugas:
1. Analisis pola dan tren dari aktivitas user
2. Identifikasi area yang bisa ditingkatkan
3. Rekomendasikan 2-3 habit baru yang spesifik dan realistis
4. Pastikan tidak duplikasi dengan habit yang sudah ada
5. Berikan alasan kenapa habit ini cocok berdasarkan pola aktivitas
6. Sertakan tips implementasi yang praktis

Format:
- Gunakan bahasa Indonesia yang friendly dan motivational
- Berikan rekomendasi dalam format yang actionable
- Fokus pada habit yang sustainable dan achievable
`, existingHabitsStr.String(), activityAnalysis.String())

	return g.callGeminiAPI(prompt)
}

// GenerateInsights generates AI-powered user insights
func (g *GeminiService) GenerateInsights(insights *models.UserInsights) (string, error) {
	prompt := fmt.Sprintf(`
Analisis data produktivitas user dan berikan insights yang valuable:

Data User:
- Total aktivitas: %d
- Total waktu: %.1f jam
- Total pengeluaran: Rp%d
- Habit aktif: %d
- Rata-rata jam harian: %.1f jam
- Waktu paling produktif: %s
- Jenis aktivitas utama: %s
- Pola pengeluaran: %s

Tugas:
1. Berikan 3-4 insights penting tentang pola produktivitas
2. Identifikasi kekuatan dan area improvement
3. Berikan rekomendasi actionable untuk optimalisasi
4. Gunakan data untuk memberikan perspektif yang personal

Format: Paragraf yang mengalir, bahasa Indonesia yang engaging dan insightful.
`,
		insights.TotalActivities,
		insights.TotalHours,
		insights.TotalExpenses,
		insights.ActiveHabits,
		insights.AvgDailyHours,
		insights.MostProductiveTime,
		insights.TopActivityType,
		insights.SpendingPattern,
	)

	return g.callGeminiAPI(prompt)
}

// AnalyzeActivities generates activity analysis
func (g *GeminiService) AnalyzeActivities(activities []models.Activity, days int) (string, error) {
	// Calculate metrics
	totalTime := 0
	totalCost := 0
	activityTypes := make(map[string]int)
	dailyDistribution := make(map[string]int)

	for _, activity := range activities {
		totalTime += activity.DurationMins
		if activity.Cost != nil {
			totalCost += *activity.Cost
		}

		category := categorizeActivity(activity.Title)
		activityTypes[category]++

		day := activity.StartTime.Weekday().String()
		dailyDistribution[day]++
	}

	avgDailyTime := float64(totalTime) / float64(days) / 60.0

	prompt := fmt.Sprintf(`
Analisis pola aktivitas user dalam %d hari terakhir:

Statistik:
- Total aktivitas: %d
- Total waktu: %d menit (%.1f jam)
- Rata-rata harian: %.1f jam
- Total pengeluaran: Rp%d
- Rata-rata per aktivitas: %.1f menit

Tugas:
1. Analisis pola dan tren produktivitas
2. Identifikasi peak performance times
3. Evaluasi efisiensi penggunaan waktu
4. Berikan rekomendasi improvement yang spesifik
5. Highlight pencapaian positif

Format: Analisis yang comprehensive tapi easy to digest, bahasa Indonesia yang professional namun friendly.
`,
		days,
		len(activities),
		totalTime,
		float64(totalTime)/60.0,
		avgDailyTime,
		totalCost,
		float64(totalTime)/float64(len(activities)),
	)

	return g.callGeminiAPI(prompt)
}

// GenerateProductivityTips generates personalized productivity tips
func (g *GeminiService) GenerateProductivityTips(context *models.UserContext) (string, error) {
	prompt := fmt.Sprintf(`
Berikan tips produktivitas yang personal untuk user dengan profil:

User: %s
- Total aktivitas: %d
- Total habits: %d
- Rata-rata jam harian: %.1f jam
- Pola recent: %s

Tugas:
1. Berikan 5-7 tips produktivitas yang personal dan actionable
2. Base on pola aktivitas dan level engagement user
3. Mix antara tips immediate action dan long-term strategy
4. Sesuaikan dengan current productivity level
5. Berikan tips yang realistis dan achievable

Format: 
- Numbered list dengan penjelasan singkat
- Bahasa Indonesia yang motivational
- Fokus pada practical implementation
`,
		context.Username,
		context.TotalActivities,
		context.TotalHabits,
		context.AvgDailyHours,
		context.RecentPatterns,
	)

	return g.callGeminiAPI(prompt)
}

// callGeminiAPI makes actual API call to Gemini
func (g *GeminiService) callGeminiAPI(prompt string) (string, error) {
	// Construct request
	request := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	// Marshal request
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s?key=%s", g.baseURL, g.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// categorizeActivity categorizes activity based on title (simplified)
func categorizeActivity(title string) string {
	title = strings.ToLower(title)

	if strings.Contains(title, "belajar") || strings.Contains(title, "study") ||
		strings.Contains(title, "learning") || strings.Contains(title, "ngoding") ||
		strings.Contains(title, "coding") || strings.Contains(title, "programming") {
		return "Learning & Development"
	}

	if strings.Contains(title, "makan") || strings.Contains(title, "food") ||
		strings.Contains(title, "breakfast") || strings.Contains(title, "lunch") ||
		strings.Contains(title, "dinner") {
		return "Food & Nutrition"
	}

	if strings.Contains(title, "olahraga") || strings.Contains(title, "exercise") ||
		strings.Contains(title, "gym") || strings.Contains(title, "workout") {
		return "Health & Fitness"
	}

	if strings.Contains(title, "meeting") || strings.Contains(title, "work") ||
		strings.Contains(title, "project") || strings.Contains(title, "task") {
		return "Work & Professional"
	}

	if strings.Contains(title, "test") || strings.Contains(title, "testing") {
		return "Testing & Quality Assurance"
	}

	return "General Activities"
}
