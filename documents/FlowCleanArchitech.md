# Clean Architecture — Phân Tích & Hướng Dẫn Toàn Diện

> _Tài liệu này được viết bởi system review và fix tự động. Mục đích: giải thích toàn bộ triết lý tổ chức code theo Clean Architecture trong hệ thống này — từ cấu trúc layers, cách tạo một API mới từ đầu đến khi có endpoint, cách tích hợp bên thứ 3, đến nơi đặt helper và constant._

---

## 1. Triết Lý Cốt Lõi (Core Philosophy)

Clean Architecture đặt ra một quy tắc bất biến:

> **Dependency Rule**: Mọi dependency đều phải trỏ vào bên trong (hướng về domain). Layer bên ngoài có thể biết về layer bên trong, nhưng KHÔNG BAO GIỜ ngược lại.

```
┌─────────────────────────────────────────────────────────────┐
│  Delivery Layer (HTTP Handlers, Queue Jobs, WebSocket)      │  ← Biết mọi thứ
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Infrastructure Layer (Repository Impl, Adapters)    │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │  Usecase Layer (Application Business Logic)    │ │  │
│  │  │  ┌───────────────────────────────────────────┐ │ │  │
│  │  │  │  Domain Layer (Entities, Interfaces)     │ │ │  │
│  │  │  │  (không import bất cứ gì từ các layer   │ │ │  │
│  │  │  │   khác trong internal/)                  │ │ │  │
│  │  │  └───────────────────────────────────────────┘ │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                    Dependencies flow INWARD ←
```

---

## 2. Cấu Trúc Thư Mục Trong Hệ Thống Này

```
internal/
├── domain/                         # LAYER 1 — Innermost. Thuần Go, không import infra.
│   ├── apperror/                   # Kiểu lỗi thuần domain (DomainError)
│   ├── consts/                     # Hằng số nghiệp vụ (error codes, queue types, status)
│   ├── model/                      # Entities + DTOs (structs đại diện bảng / request / response)
│   │   ├── user/                   # Model theo domain slice (user.go, upload_avatar.go...)
│   │   ├── chatgroup/
│   │   └── ...
│   ├── repository/                 # Interfaces cho data access (không có implementation)
│   ├── service/                    # Interfaces cho external services (Cache, File, Queue, OAuth...)
│   └── usecase/                    # Interfaces cho application logic
│
├── usecase/                        # LAYER 2 — Application business rules
│   ├── user/user.go                # Implement domain/usecase.UserUsecase
│   ├── auth/auth.go                # Implement domain/usecase.AuthUsecase
│   ├── auth/oauth.go               # Implement Google OAuth flow
│   ├── totp/totp.go                # Implement domain/usecase.TotpUsecase
│   ├── articles.go                 # Implement domain/usecase.ArticlesUsecase
│   └── ...
│
├── repository/                     # LAYER 3 — Data access implementations
│   ├── user.go                     # Implement domain/repository.UserRepository (GORM)
│   ├── ai.go                       # Implement domain/repository.AiRepository (HTTP call)
│   └── ...
│
├── infrastructure/                 # LAYER 3 — Infrastructure adapters
│   └── fileadapter/adapter.go      # Adapter: pkg/file.HTTPFile → domain/service.FileService
│
├── server/                         # LAYER 4 — Delivery (ngoài cùng)
│   ├── http/handler/               # HTTP handlers (Echo)
│   │   ├── user/handler.go         # Route registration
│   │   ├── user/user.go            # HTTP handler functions
│   │   └── ...
│   └── queue/jobs/                 # Queue job handlers
│       ├── mail_job.go
│       └── ...
│
├── config/                         # App config structs + initializer
└── di/wire.go                      # Composition Root — nơi duy nhất biết tất cả mọi thứ

pkg/                                # Shared infrastructure utilities (reusable, no business logic)
├── cache/                          # Redis implementation
├── http/                           # Echo helpers, error types, response functions
├── logger/                         # Zerolog wrapper
├── mailer/                         # SMTP mailer
├── queue/                          # In-memory queue
├── file/                           # File download/upload utilities
├── websocket/                      # Gorilla WebSocket hub + client
└── ...
```

---

## 3. Các Quy Tắc Import (Dependency Rules)

| Layer               | Được phép import                                      | KHÔNG được import                                                        |
| ------------------- | ----------------------------------------------------- | ------------------------------------------------------------------------ |
| **domain/**         | `pkg/query`, stdlib                                   | `internal/usecase`, `internal/repository`, `internal/server`, `pkg/http` |
| **usecase/**        | `domain/*`, `pkg/logger`, stdlib                      | `internal/repository` (impl), `internal/server`, `pkg/http` (errors)     |
| **repository/**     | `domain/*`, `pkg/logger`, `pkg/http` (client), stdlib | `internal/usecase`, `internal/server`                                    |
| **infrastructure/** | `domain/*`, `pkg/*`                                   | `internal/usecase`, `internal/server`                                    |
| **server/**         | tất cả                                                | — (đây là layer ngoài cùng)                                              |
| **di/wire.go**      | tất cả                                                | — (đây là Composition Root)                                              |

### ⚠️ Lỗi Phổ Biến Cần Tránh

```go
// ❌ SAI: Domain import pkg/http (HTTP là delivery concern)
import xhttp "thomas.vn/.../pkg/http"
func EmailAlreadyExistsError() *xhttp.AppError { ... }

// ✅ ĐÚNG: Domain dùng apperror (pure Go, no framework)
import "thomas.vn/.../internal/domain/apperror"
func EmailAlreadyExistsError() *apperror.DomainError { ... }

// ❌ SAI: Usecase import pkg/http để tạo error
return nil, xhttp.BadRequestErrorf("invalid input")

// ✅ ĐÚNG: Usecase dùng domain apperror
return nil, apperror.BadRequest("invalid input")

// ❌ SAI: Usecase dùng concrete struct từ repository implementation
aiRepo repository.AiRepository   // struct, không phải interface

// ✅ ĐÚNG: Usecase dùng interface từ domain/repository
aiRepo domainrepo.AiRepository   // interface → có thể mock trong test
```

---

## 4. Tạo Một API Mới (End-to-End Flow)

Ví dụ: Tạo API `POST /api/apartments` để đăng apartment mới.

### Bước 1 — Định nghĩa Entity trong Domain

```go
// internal/domain/model/apartment/apartment.go
package apartment

import "time"

// Apartment là entity nghiệp vụ, ánh xạ trực tiếp với bảng `apartments`
type Apartment struct {
    ID          int       `json:"id" gorm:"primary_key"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    UserID      int       `json:"user_id"`
    IsDeleted   bool      `json:"is_deleted"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func (Apartment) TableName() string { return "apartments" }

// CreateApartmentRequest là DTO nhận từ HTTP request
type CreateApartmentRequest struct {
    Title       string  `json:"title"       validate:"required,min=5"`
    Description string  `json:"description" validate:"required"`
    Price       float64 `json:"price"       validate:"required,gt=0"`
    UserID      int     `json:"-"` // lấy từ JWT, không từ body
}

// ApartmentIDRequest dùng cho các route có path param :id
type ApartmentIDRequest struct {
    ID uint `param:"id" validate:"required,gt=0"`
}
```

### Bước 2 — Định nghĩa Repository Interface trong Domain

```go
// internal/domain/repository/apartment.go
package repository

import (
    "context"
    apt "thomas.vn/.../internal/domain/model/apartment"
)

// ApartmentRepository định nghĩa contract truy cập dữ liệu.
// Usecase layer PHỤ THUỘC vào interface này, không phải implementation.
type ApartmentRepository interface {
    Create(ctx context.Context, a *apt.Apartment) (*apt.Apartment, error)
    GetByID(ctx context.Context, id uint) (*apt.Apartment, error)
    ListByUser(ctx context.Context, userID int) ([]*apt.Apartment, error)
    Update(ctx context.Context, a *apt.Apartment) (*apt.Apartment, error)
    Delete(ctx context.Context, id uint) error
}
```

### Bước 3 — Định nghĩa Usecase Interface trong Domain

```go
// internal/domain/usecase/apartment.go
package usecase

import (
    "context"
    apt "thomas.vn/.../internal/domain/model/apartment"
)

// ApartmentUsecase định nghĩa application logic contract.
// Handler phụ thuộc vào interface này.
type ApartmentUsecase interface {
    Create(ctx context.Context, req *apt.CreateApartmentRequest) (*apt.Apartment, error)
    GetByID(ctx context.Context, id uint) (*apt.Apartment, error)
    ListByUser(ctx context.Context, userID int) ([]*apt.Apartment, error)
    Update(ctx context.Context, id uint, req *apt.CreateApartmentRequest) (*apt.Apartment, error)
    Delete(ctx context.Context, id uint) error
}
```

### Bước 4 — Implement Repository (Infrastructure Layer)

```go
// internal/repository/apartment.go
package repository

import (
    "errors"
    "gorm.io/gorm"
    "thomas.vn/.../internal/domain/apperror"
    domainrepo "thomas.vn/.../internal/domain/repository"
    apt "thomas.vn/.../internal/domain/model/apartment"
    xlogger "thomas.vn/.../pkg/logger"
)

type apartmentRepository struct {
    logger *xlogger.Logger
    db     *gorm.DB
}

// Đảm bảo implementation thỏa mãn interface tại compile time
var _ domainrepo.ApartmentRepository = (*apartmentRepository)(nil)

func NewApartmentRepository(logger *xlogger.Logger, db *gorm.DB) domainrepo.ApartmentRepository {
    return &apartmentRepository{logger: logger, db: db.Table("apartments")}
}

func (r *apartmentRepository) Create(ctx context.Context, a *apt.Apartment) (*apt.Apartment, error) {
    if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
        r.logger.Error("Create apartment failed", xlogger.Error(err))
        return nil, err
    }
    return a, nil
}

func (r *apartmentRepository) GetByID(ctx context.Context, id uint) (*apt.Apartment, error) {
    var a apt.Apartment
    err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = false", id).First(&a).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil // không tìm thấy trả nil, usecase xử lý
    }
    return &a, err
}

// ... các method khác tương tự
```

### Bước 5 — Implement Usecase (Application Layer)

```go
// internal/usecase/apartment/apartment.go
package apartment

import (
    "context"
    "thomas.vn/.../internal/domain/apperror"
    domainrepo "thomas.vn/.../internal/domain/repository"
    domainuc "thomas.vn/.../internal/domain/usecase"
    apt "thomas.vn/.../internal/domain/model/apartment"
    xlogger "thomas.vn/.../pkg/logger"
)

type apartmentUsecase struct {
    logger *xlogger.Logger
    repo   domainrepo.ApartmentRepository
}

var _ domainuc.ApartmentUsecase = (*apartmentUsecase)(nil)

func NewApartmentUsecase(logger *xlogger.Logger, repo domainrepo.ApartmentRepository) domainuc.ApartmentUsecase {
    return &apartmentUsecase{logger: logger, repo: repo}
}

func (u *apartmentUsecase) Create(ctx context.Context, req *apt.CreateApartmentRequest) (*apt.Apartment, error) {
    entity := &apt.Apartment{
        Title:       req.Title,
        Description: req.Description,
        Price:       req.Price,
        UserID:      req.UserID,
    }
    return u.repo.Create(ctx, entity)
}

func (u *apartmentUsecase) GetByID(ctx context.Context, id uint) (*apt.Apartment, error) {
    a, err := u.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if a == nil {
        // ✅ Trả DomainError, không phải HTTP error
        return nil, apperror.NotFound("Apartment with ID %d not found", id)
    }
    return a, nil
}
// ... các method khác
```

### Bước 6 — Tạo HTTP Handler (Delivery Layer)

```go
// internal/server/http/handler/apartment/handler.go
package apartment

import (
    "thomas.vn/.../internal/domain/usecase"
    xlogger "thomas.vn/.../pkg/logger"
)

// Handler là entry point cho tất cả routes của apartment.
// Dùng Functional Options Pattern để thuận tiện test và mở rộng.
type Handler struct {
    logger           *xlogger.Logger
    apartmentHandler *ApartmentHandler
}

type HandlerOption func(*Handler)

func WithApartmentUsecase(uc usecase.ApartmentUsecase) HandlerOption {
    return func(h *Handler) {
        h.apartmentHandler = NewApartmentHandler(h.logger, uc)
    }
}

func NewHandler(logger *xlogger.Logger, opts ...HandlerOption) *Handler {
    h := &Handler{logger: logger}
    for _, opt := range opts { opt(h) }
    return h
}

func (h *Handler) Apartment() *ApartmentHandler { return h.apartmentHandler }
```

```go
// internal/server/http/handler/apartment/apartment.go
package apartment

import (
    "github.com/labstack/echo/v4"
    domainuc "thomas.vn/.../internal/domain/usecase"
    apt "thomas.vn/.../internal/domain/model/apartment"
    xhttp "thomas.vn/.../pkg/http"
    xcontext "thomas.vn/.../pkg/http/context"
    xlogger "thomas.vn/.../pkg/logger"
)

type ApartmentHandler struct {
    logger     *xlogger.Logger
    apartmentUC domainuc.ApartmentUsecase
}

func NewApartmentHandler(logger *xlogger.Logger, uc domainuc.ApartmentUsecase) *ApartmentHandler {
    return &ApartmentHandler{logger: logger, apartmentUC: uc}
}

// Create godoc
// @Summary Create apartment
// @Tags apartments
// @Accept json
// @Produce json
// @Param data body apt.CreateApartmentRequest true "Create apartment"
// @Success 201 {object} xhttp.APIResponse{data=apt.Apartment}
// @Router /api/apartments [post]
func (h *ApartmentHandler) Create(c echo.Context) error {
    var req apt.CreateApartmentRequest
    if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
        return xhttp.BadRequestResponse(c, err)
    }

    // Lấy User ID từ JWT context (set bởi Auth middleware)
    req.UserID = int(xcontext.GetUserID(c))

    res, err := h.apartmentUC.Create(c.Request().Context(), &req)
    if err != nil {
        h.logger.Error("Create apartment failed", xlogger.Error(err))
        return xhttp.AppErrorResponse(c, err) // tự động handle DomainError + AppError
    }

    return xhttp.CreatedResponse(c, res)
}
```

### Bước 7 — Đăng ký Route

```go
// internal/server/http/handler/root/ (hoặc nơi đăng ký route)
// Thêm apartment routes vào router chính:

func registerApartmentRoutes(e *echo.Echo, h *apartment.Handler, authMw echo.MiddlewareFunc) {
    api := e.Group("/api", authMw)
    api.POST("/apartments", h.Apartment().Create)
    api.GET("/apartments/:id", h.Apartment().GetByID)
    api.PUT("/apartments/:id", h.Apartment().Update)
    api.DELETE("/apartments/:id", h.Apartment().Delete)
}
```

### Bước 8 — Kết nối trong DI (Composition Root)

```go
// internal/di/wire.go
// Thêm vào NewAppContainer:

// Repository
apartmentRepo := repository.NewApartmentRepository(logger, mysqlClient.DB)

// Usecase
apartmentUC := apt_uc.NewApartmentUsecase(logger, apartmentRepo)

// Handler
apartmentHandler := apartment_handler.NewHandler(logger,
    apartment_handler.WithApartmentUsecase(apartmentUC),
)

// Register routes
registerApartmentRoutes(httpServer, apartmentHandler, authMiddleware)
```

---

## 5. Tích Hợp Bên Thứ 3 (Third-Party Integration)

Khi tích hợp dịch vụ bên ngoài (OAuth, Payment, SMS, AI API...), áp dụng nguyên tắc:

> **Domain định nghĩa interface → Infrastructure implement → Adapter bridge nếu cần**

### 5.1 External HTTP API (VD: AI Service, Payment Gateway)

```
internal/
├── domain/
│   ├── repository/ai.go        ← Interface: AiRepository { VerifyCV(...) }
│   └── model/ai.go             ← Request/Response DTOs thuần Go
├── repository/
│   └── ai.go                   ← Implementation: gọi HTTP API bên ngoài
└── di/wire.go                  ← Wiring: NewAiRepository(httpClient, fileSvc, url)
```

```go
// ✅ Domain chỉ biết interface
type AiRepository interface {
    VerifyCV(attachFile, jobDesc string) (int, model.VerifyResponse, error)
}

// ✅ Implementation ở infrastructure layer
type AiRepository struct {
    client  *xhttp.HTTPClient  // gọi HTTP
    fileSvc service.FileService // download file
}
```

### 5.2 OAuth Provider (Google, Facebook)

```
internal/
├── domain/
│   └── service/oauth.go      ← Interface: GoogleOAuthService { GetProfile(...) }
└── di/wire.go                ← Wired với xgoogle.New(clientID, secret, callbackURL)

pkg/
└── oauth/
    └── google/               ← Implementation, không có business logic
```

### 5.3 File Storage (Local / Cloudinary)

```
internal/
├── domain/
│   └── service/file.go       ← Interface: FileService { Upload, Download, Delete... }
│                                          File struct (pure Go)
├── infrastructure/
│   └── fileadapter/          ← Adapter: pkg/file.HTTPFile → service.FileService
└── di/wire.go                ← fileImpl := xfile.NewHTTPFile(...)
                                 fileSvc  := fileadapter.New(fileImpl)
```

> **Tại sao cần Adapter?**
> `pkg/file.HTTPFile.Download()` trả về `xfile.File` (infrastructure type).
> `domain/service.FileService.Download()` trả về `service.File` (domain type).
> Adapter ở `internal/infrastructure/fileadapter/` convert giữa 2 types.
> → **Domain hoàn toàn không biết đến pkg/file.**

### 5.4 Email Service

```
internal/
├── config/mailer.go          ← MailerConfig { SMTP {Host,Port,User,Pass}, FromName }
├── config/config.go          ← Config { ..., Mailer MailerConfig }
├── domain/
│   └── repository/mail.go   ← Interface: MailRepository { Send(ctx, MailData) error }
└── di/wire.go                ← mailer := mail.NewMailer(cfg.Mailer.SMTP, cfg.Mailer.From)

pkg/
└── mailer/                   ← SMTP implementation. Credentials từ Config, KHÔNG hardcode.
```

**⚠️ Security rule: KHÔNG BAO GIỜ hardcode credentials trong source code.**

```go
// ❌ SAI — credentials trong code
mailer := mail.NewMailer(mail.SMTPConfig{
    User: "abc@gmail.com",
    Pass: "my_secret_password",
}, "My App")

// ✅ ĐÚNG — credentials từ config (loaded từ YAML / env vars)
mailer := mail.NewMailer(mail.SMTPConfig{
    Host: cfg.Mailer.SMTP.Host,
    Port: cfg.Mailer.SMTP.Port,
    User: cfg.Mailer.SMTP.User,
    Pass: cfg.Mailer.SMTP.Pass,
}, cfg.Mailer.FromName+" <"+cfg.Mailer.SMTP.User+">")
```

### 5.5 Message Queue (In-Memory / Redis / RabbitMQ)

```
internal/
├── domain/
│   ├── consts/queue.go       ← MessageType constants (type alias = string, no pkg/queue import)
│   └── service/queue.go      ← Interface: QueueService { PublishMessage(ctx, string, payload) error }
└── server/
    └── queue/jobs/           ← Job implementations (Delivery layer, biết về xqueue)
```

```go
// ✅ Domain consts không import pkg/queue
type MessageType = string  // type alias

const (
    MailJobType       MessageType = "mail_job"
    QueueMailRegister MessageType = "mail_register"
)

// ✅ Domain service interface dùng string (tương thích với bất kỳ queue implementation nào)
type QueueService interface {
    PublishMessage(ctx context.Context, msgType string, payload interface{}) error
}
```

---

## 6. Domain Errors — Tầng Nào Biết Gì

```
apperror.DomainError (internal/domain/apperror/)
         ↑ Usecase tạo ra, trả về
         ↑ Handler nhận, convert sang HTTP response
         ↑ pkg/http/response.go detect via interface (không import internal)
```

```go
// internal/domain/apperror/apperror.go
type DomainError struct {
    Code    string // machine-readable: "ERR_NOT_FOUND"
    Message string // human-readable:  "User not found"
    Field   string // field gây lỗi (optional): "email"
    Status  int    // HTTP hint: 404
}

// Methods implemented → pkg/http nhận diện qua so sánh interface cục bộ:
func (e *DomainError) GetStatus() int   { return e.Status }
func (e *DomainError) GetCode() string  { return e.Code }
func (e *DomainError) GetField() string { return e.Field }
```

```go
// pkg/http/response.go
// Sử dụng local interface để nhận diện DomainError mà KHÔNG cần import internal/
type domainError interface {
    error
    GetStatus() int
    GetCode() string
    GetField() string
}

func AppErrorResponse(c echo.Context, err error) error {
    // 1. HTTP-layer errors (AppError)
    var appErr *AppError
    if errors.As(err, &appErr) { ... }

    // 2. Domain errors (DomainError) → duck typing, no import needed
    var de domainError
    if errors.As(err, &de) {
        return c.JSON(de.GetStatus(), APIResponse{...})
    }

    return InternalServerErrorResponse(c) // fallback
}
```

---

## 7. Constants và Helpers — Đặt Ở Đâu

| Loại                      | Đặt ở đâu                          | Ví dụ                                            |
| ------------------------- | ---------------------------------- | ------------------------------------------------ |
| Domain error codes        | `internal/domain/consts/errors.go` | `ErrEmailAlreadyExists = "ERR_..."`              |
| Domain error constructors | `internal/domain/consts/errors.go` | `func EmailAlreadyExistsError(...) *DomainError` |
| Queue message types       | `internal/domain/consts/queue.go`  | `MailJobType MessageType = "mail_job"`           |
| User status constants     | `internal/domain/consts/user.go`   | `UserStatusActive = 1`                           |
| Business logic helpers    | `internal/domain/model/<feature>/` | helper methods trên struct                       |
| Generic string/time utils | `pkg/utils/`                       | `string.go`, `time.go`                           |
| HTTP form helpers         | `pkg/http/form.go`                 | parse request, validate                          |
| Pagination/Sort/DateRange | `pkg/query/`                       | `PaginationOptions`, `SortOptions`               |
| File constants            | `pkg/file/consts.go`               | allowed mime types                               |
| Crypto/hash utils         | `pkg/utils/`                       | hash functions                                   |

**Rule of thumb:**

- Liên quan đến nghiệp vụ → `internal/domain/`
- Liên quan đến kỹ thuật thuần (không có business logic) → `pkg/`
- Liên quan đến HTTP delivery → `pkg/http/` hoặc `internal/server/http/`
- Liên quan đến database query → `pkg/query/` (pagination, sort, date range)

---

## 8. Tổng Hợp Các Vi Phạm Đã Sửa (Changelog)

### ❌ → ✅ Violations Fixed

| File                       | Vi phạm                                                                       | Đã sửa                                            |
| -------------------------- | ----------------------------------------------------------------------------- | ------------------------------------------------- |
| `domain/consts/errors.go`  | Import `pkg/http` → Domain phụ thuộc HTTP                                     | Dùng `domain/apperror`                            |
| `domain/consts/queue.go`   | Import `pkg/queue` → Domain phụ thuộc infra queue                             | Dùng native `type MessageType = string`           |
| `domain/service/cache.go`  | Embedding `xcache.CacheService` → lộ infra type                               | Redeclare interface explicitly                    |
| `domain/service/file.go`   | Embedding `xfile.FileService` → lộ infra type                                 | Domain `File` struct + explicit interface         |
| `domain/service/queue.go`  | Embedding `xqueue.QueueService` → lộ infra type                               | Explicit interface với `string` msgType           |
| `usecase/user/user.go`     | `xhttp.BadRequestErrorf` + `xqueue.QueueService` trong usecase                | Dùng `apperror` + `service.QueueService`          |
| `usecase/auth/auth.go`     | Same violations as above                                                      | Same fixes                                        |
| `usecase/auth/oauth.go`    | `xhttp.BadRequestErrorf` trong usecase                                        | Dùng `apperror`                                   |
| `usecase/chatgroup.go`     | `xhttp.BadRequestErrorf` trong usecase                                        | Dùng `apperror`                                   |
| `usecase/permission.go`    | `xhttp.NotFoundErrorf` trong usecase                                          | Dùng `apperror`                                   |
| `usecase/totp/totp.go`     | `xhttp.BadRequestErrorf` trong usecase                                        | Dùng `apperror`                                   |
| `usecase/ai.go`            | Dùng concrete struct `repository.AiRepository`                                | Dùng `domain/repository.AiRepository` interface   |
| `di/wire.go`               | Hardcoded SMTP email + password                                               | Đọc từ `cfg.Mailer`                               |
| `pkg/mailer/smtp.go`       | Hardcoded sender email string                                                 | Dùng `m.senderEmail` field                        |
| `pkg/websocket/server.go`  | Dùng `golang.org/x/net/websocket` (cũ) + exported Client fields không tồn tại | Rewrite dùng gorilla/websocket, unexported fields |
| `pkg/websocket/hub.go`     | `c.Send` (exported) không tồn tại                                             | Dùng `c.send` (unexported, same package)          |
| `pkg/websocket/handler.go` | Dùng `golang.org/x/net/websocket` handler cũ                                  | Gorilla upgrader + Echo pattern                   |

### 📁 Files Mới

| File                                             | Mục đích                                     |
| ------------------------------------------------ | -------------------------------------------- |
| `internal/domain/apperror/apperror.go`           | Pure domain error type                       |
| `internal/domain/repository/ai.go`               | AiRepository interface (previously missing)  |
| `internal/config/mailer.go`                      | MailerConfig struct                          |
| `internal/infrastructure/fileadapter/adapter.go` | Bridge pkg/file ↔ domain/service.FileService |

---

## 9. Checklist Khi Tạo Feature Mới

```
□ 1. Tạo model (entity + request DTOs) trong internal/domain/model/<feature>/
□ 2. Tạo repository interface trong internal/domain/repository/<feature>.go
□ 3. Tạo usecase interface trong internal/domain/usecase/<feature>.go
□ 4. Implement repository trong internal/repository/<feature>.go
□ 5. Implement usecase trong internal/usecase/<feature>/<feature>.go
□ 6. Tạo HTTP handler trong internal/server/http/handler/<feature>/
□ 7. Đăng ký route trong root handler
□ 8. Wire tất cả trong internal/di/wire.go
□ 9. Thêm config vào config/base.yaml nếu cần
□ 10. Viết unit test cho usecase (mock repository)

Quy tắc kiểm tra:
□ Domain layer không import pkg/http, pkg/queue (trừ pkg/query cho pagination)
□ Usecase không import pkg/http (dùng apperror thay thế)
□ Usecase không dùng concrete struct từ repository (dùng domain interface)
□ Không có credentials hardcode trong source code
□ Mọi external service đều có interface trong domain/service/ hoặc domain/repository/
```
