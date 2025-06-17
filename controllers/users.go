package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)


type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}


type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Name            string `json:"name"`
}

func RegisterUser(e *core.RequestEvent) error {
	var req RegisterRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.PasswordConfirm == "" {
		return e.BadRequestError("Email and password required", nil)
	}
	if req.Password != req.PasswordConfirm {
		return e.BadRequestError("Passwords do not match", nil)
	}

	usersCollection, err := e.App.FindCollectionByNameOrId("users")
	if err != nil {
		return e.InternalServerError("Users collection not found", err)
	}

	// Check if user already exists
	existing, err := e.App.FindFirstRecordByData("users", "email", req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return e.InternalServerError("Error checking for existing user", err)
	}
	if existing != nil {
		return e.BadRequestError("User already exists", nil)
	}

	// Create new user record
	record := core.NewRecord(usersCollection)
	record.Load(map[string]any{
		"email":    req.Email,
		"password": req.Password,
		"name":     req.Name,
	})

	if err := e.App.Save(record); err != nil {
		return e.BadRequestError("Could not create user", err)
	}

// 	if err := e.App.SendEmailVerification(record); err != nil {
// }

	return e.JSON(http.StatusCreated, record.PublicExport())
}

func PathCheck(e *core.RequestEvent) error {
	return e.JSON(http.StatusOK, map[string]string{
		"message": "User route is working",
	})
}


func LoginUser(e *core.RequestEvent) error {
    var req LoginRequest
    if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
        return e.BadRequestError("Invalid request body", err)
    }

    if req.Email == "" || req.Password == "" {
        return e.BadRequestError("Email and password required", nil)
    }

    // Find user by email
    record, err := e.App.FindFirstRecordByData("users", "email", req.Email)
    if err != nil {
        return e.BadRequestError("Invalid email or password", nil)
    }

    // Check password (assuming plain text for example; use hashing in production)
    if record.GetString("password") != req.Password {
        return e.BadRequestError("Invalid email or password", nil)
    }

    // // Optionally check if email is verified
    // if !record.GetBool("verified") {
    //     return e.BadRequestError("Email not verified", nil)
    // }

    return e.JSON(http.StatusOK, map[string]any{
        "message": "Login successful",
        "user":    record.PublicExport(),
    })
}