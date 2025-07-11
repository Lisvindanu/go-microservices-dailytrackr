# 🧠 DailyTrackr Microservices

# 🧠 Go Microservices DailyTrackr

A productivity tracking platform built with **Go microservices architecture**, featuring habit tracking, activity logging, and AI-powered insights.

## 🚧 Development Status

**Overall Progress: 35% Complete**

### ✅ Completed Features
- [x] **Microservices Foundation** (100%)
    - [x] Shared package with DTOs, utils, config
    - [x] Database schema and MySQL integration
    - [x] Environment configuration system
- [x] **User Service** (100%)
    - [x] JWT-based authentication
    - [x] User registration and login
    - [x] Password encryption with bcrypt
    - [x] User profile management
- [x] **Gateway Service** (100%)
    - [x] API Gateway with request routing
    - [x] Service proxy functionality
    - [x] CORS middleware
    - [x] Health check endpoints

### 🔄 In Development
- [ ] **Activity Service** (0%)
    - [ ] Activity CRUD operations
    - [ ] Photo upload integration
    - [ ] Time and cost tracking
- [ ] **Habit Service** (0%)
    - [ ] 30-day habit challenges
    - [ ] Daily habit logging
    - [ ] Progress tracking and streaks
- [ ] **AI Service** (0%)
    - [ ] Gemini API integration
    - [ ] Daily activity summaries
    - [ ] Habit recommendations

### 📋 Planned Features
- [ ] **Notification Service** (0%)
    - [ ] Email reminders
    - [ ] Daily and weekly reports
    - [ ] Habit reminder notifications
- [ ] **Statistics Service** (0%)
    - [ ] Activity analytics
    - [ ] Habit success rates
    - [ ] Progress visualizations
- [ ] **Frontend Application** (0%)
    - [ ] React.js web application
    - [ ] User dashboard
    - [ ] Activity and habit management UI

## 🏗️ Current Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Gateway   │    │ User Service│    │Activity Svc │
│   :3000 ✅  │────│   :3001 ✅  │    │  :3002 🔄   │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                   ┌─────────────┐    ┌─────────────┐
                   │Habit Service│    │  AI Service │
                   │  :3003 🔄   │    │  :3006 🔄   │
                   └─────────────┘    └─────────────┘
                           │
                   ┌─────────────────┐
                   │   MySQL DB ✅   │
                   │   dailytrackr   │
                   └─────────────────┘
```

**Legend:** ✅ Complete | 🔄 In Development | ⏳ Planned

## 🔧 Development Progress

### Service Status
| Service | Port | Status | Progress |
|---------|------|--------|----------|
| Gateway | 3000 | ✅ Working | 100% |
| User Service | 3001 | ✅ Working | 100% |
| Activity Service | 3002 | 🔄 Development | 0% |
| Habit Service | 3003 | 🔄 Development | 0% |
| Notification Service | 3004 | ⏳ Planned | 0% |
| Statistics Service | 3005 | ⏳ Planned | 0% |
| AI Service | 3006 | 🔄 Development | 0% |

### Database Status
- ✅ **Schema Design** - 6 tables created
- ✅ **MySQL Integration** - Connection established
- ✅ **Sample Data** - Test users and data inserted
- ✅ **Migrations** - Database structure ready

## 📡 Working API Endpoints

### Gateway (Port 3000)
- ✅ `GET /` - Gateway health & service status
- ✅ `POST /api/users/auth/register` - User registration
- ✅ `POST /api/users/auth/login` - User login
- ✅ `GET /api/users/health` - User service health

### User Service (Port 3001)
- ✅ `GET /health` - Service health check
- ✅ `POST /auth/register` - User registration
- ✅ `POST /auth/login` - User authentication
- ✅ `GET /api/v1/users/profile` - Get user profile (JWT required)
- ✅ `GET /api/v1/users/:id` - Get user by ID (JWT required)

### Planned Endpoints
- 🔄 Activity Service endpoints
- 🔄 Habit Service endpoints
- 🔄 AI Service endpoints
- 🔄 Statistics endpoints

## 🛠️ Tech Stack

- **Language**: Go 1.24
- **Framework**: Gin (HTTP router)
- **Database**: MySQL 8.0
- **Architecture**: Microservices with API Gateway
- **Authentication**: JWT tokens with bcrypt
- **Validation**: go-playground/validator
- **Configuration**: godotenv
- **AI Integration**: Google Gemini API (planned)

## 📁 Project Structure

```
go-microservices-dailytrackr/
├── shared/                 # ✅ Shared utilities & DTOs
│   ├── config/            # ✅ Environment configuration
│   ├── constants/         # ✅ Application constants
│   ├── database/          # ✅ Database connection
│   ├── dto/               # ✅ Data transfer objects
│   └── utils/             # ✅ Common utilities
├── gateway/               # ✅ API Gateway service
├── user-service/          # ✅ Authentication service  
├── activity-service/      # 🔄 Activity tracking (planned)
├── habit-service/         # 🔄 Habit management (planned)
├── ai-service/           # 🔄 AI insights (planned)
├── notification-service/  # ⏳ Email notifications (planned)
├── stat-service/         # ⏳ Statistics (planned)
├── .env                  # ✅ Environment variables
└── README.md             # ✅ Documentation
```

## 🎯 Next Development Milestones

1. **Activity Service** (Target: Week 1)
    - CRUD operations for daily activities
    - Photo upload integration with Cloudinary
    - Time tracking and cost logging

2. **Habit Service** (Target: Week 2)
    - 30-day habit challenge system
    - Daily habit logging with status tracking
    - Streak calculations and progress metrics

3. **AI Service** (Target: Week 3)
    - Google Gemini API integration
    - Automated daily activity summaries
    - Smart habit recommendations

4. **Frontend Application** (Target: Week 4)
    - React.js web application
    - User authentication and dashboard
    - Activity and habit management interface

## 🤝 Contributing

1. Fork the project
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Lisvindanu**
- GitHub: [@lisvindanu](https://github.com/lisvindanu)
- Email: Lisvindanu015@gmail.com

## 🙏 Acknowledgments

- Built with ❤️ using Go microservices
- Powered by Google Gemini AI
- Database design optimized for productivity tracking

---

⭐ **Star this repo if you find it helpful!**