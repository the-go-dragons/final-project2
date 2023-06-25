# final-project2
### Final Project2 of Software Engineering Bootcamp with Go Programming Language by The Go Dragons Team
```
.
├── cmd ------------------------------------> Entry point
│   └── main.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── app --------------------------------> Initialization 
│   │   └── app.go
│   ├── domain -----------------------------> Encapsulate Business Rules (Entities Layer)
│   │   └── user.go
│   ├── interfaces -------------------------> Interface Adapters (Controller Layer)
│   │   ├── http ---------------------------> HTTP Handlers
│   │   │   ├── login_handler.go
│   │   │   ├── logout_handler.go
│   │   │   ├── signup_handler.go
│   │   │   ├── user_handler.go
│   │   │   └── utils.go
│   │   └── persistence --------------------> Provides Abstract Layer to Database
│   │       └── user_repository.go
│   └── usecase ----------------------------> The application logic (Usecase Layer)
│       └── user_usecase.go
├── LICENSE
├── pkg ------------------------------------> External Frameworks and Drivers (Port|Infra)
│   ├── config -----------------------------> Configurations
│   │   └── env.go
│   └── database ---------------------------> Databse Layer that contains GORM
│       ├── db.go
│       └── migrations ---------------------> SQL Migration files
│           └── 000001_final_project_schema.down.sql
├── README.md
└── wait-for-it.sh
```
<p style="text-align: center; width: 100%; ">Copyright&copy; 2023 <a href="https://github.com/the-go-dragons">The Go Dragons Team</a>, Licensed under MIT</p>