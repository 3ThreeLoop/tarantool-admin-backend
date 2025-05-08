package middlewares

import (
	// "api-mini-shop/internal/front/auth"
	types "api-mini-shop/pkg/share"
	"api-mini-shop/pkg/utls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	response "api-mini-shop/pkg/http/response"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/tarantool/go-tarantool/v2/pool"
)

func NewJwtMinddleWare(app *fiber.App, connPool *pool.ConnectionPool) {
	// Load environment variables
	errs := godotenv.Load()
	if errs != nil {
		log.Fatalf("Error loading .env file")
	}
	secret_key := os.Getenv("JWT_SECRET_KEY")

	// JWT Middleware
	app.Use(func(c *fiber.Ctx) error {
		// Check if the request is upgrading to WebSocket
		if websocketUpgrade := c.Get("Upgrade"); websocketUpgrade == "websocket" {
			// Extract Bearer token from Sec-WebSocket-Protocol
			webSocketProtocol := c.Get("Sec-WebSocket-Protocol")
			if webSocketProtocol == "" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Missing WebSocket protocol for authentication",
				})
			}

			// Split "Bearer, <token>"
			parts := strings.Split(webSocketProtocol, ",")
			if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid WebSocket protocol authentication format",
				})
			}

			// Extract the JWT token from the second part
			tokenString := strings.TrimSpace(parts[1])

			// Parse the JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret_key), nil
			})
			if err != nil || !token.Valid {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or expired JWT token",
				})
			}
			c.Locals("jwt_data", token)
			// Set the response header to echo back the protocol
			c.Set("Sec-WebSocket-Protocol", "Bearer")
			return c.Next()
		}

		// Apply JWT middleware for HTTP requests
		return jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret_key)},
			ContextKey: "jwt_data",
		})(c)
	})

	// Player Context Middleware
	app.Use(func(c *fiber.Ctx) error {
		// Extract the JWT token data
		player_token := c.Locals("jwt_data").(*jwt.Token)
		pclaim := player_token.Claims.(jwt.MapClaims)

		// Check if the connection is WebSocket and handle accordingly
		if websocketUpgrade := c.Get("Upgrade"); websocketUpgrade == "websocket" {
			// For WebSocket, ensure the token contains necessary claims
			return handlePlayerContext(c, pclaim, connPool)
		}

		// Handle regular HTTP requests
		return handlePlayerContext(c, pclaim, connPool)
	})
}

// Helper function to handle player context creation and session validation
func handlePlayerContext(c *fiber.Ctx, pclaim jwt.MapClaims, connPool *pool.ConnectionPool) error {
	// Check login session
	login_session, ok := pclaim["login_session"].(string)
	if !ok || login_session == "" {
		smg_error := response.NewResponseError(
			utls.Translate("loginSessionMissing", nil, c),
			-500,
			fmt.Errorf(utls.Translate("loginSessionMissing", nil, c)),
		)
		return c.Status(http.StatusUnprocessableEntity).JSON(smg_error)
	}

	// Safely extract and validate required fields
	id, ok := pclaim["id"].(float64)
	if !ok {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			"Invalid or missing 'id' in claims", -500, fmt.Errorf("missing or invalid 'id'"),
		))
	}

	playerUuid, ok := pclaim["player_uuid"].(string)
	if !ok || playerUuid == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			"Invalid or missing 'player_uuid' in claims", -500, fmt.Errorf("missing or invalid 'player_uuid'"),
		))
	}

	userName, ok := pclaim["username"].(string)
	if !ok || userName == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			"Invalid or missing 'username' in claims", -500, fmt.Errorf("missing or invalid 'username'"),
		))
	}

	membershipId, ok := pclaim["membership_id"].(float64)
	if !ok {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			"Invalid or missing 'membership_id' in claims", -500, fmt.Errorf("missing or invalid 'membership_id'"),
		))
	}

	exp, ok := pclaim["exp"].(float64)
	if !ok {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			"Invalid or missing 'exp' in claims", -500, fmt.Errorf("missing or invalid 'exp'"),
		))
	}

	status_id, ok := pclaim["status_id"].(float64)
	if !ok {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			"Invalid or missing 'status_id' in claims", -500, fmt.Errorf("missing or invalid 'status_id'"),
		))
	}

	token_version, ok := pclaim["token_version"].(float64)
	if !ok {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			"Invalid or missing 'status_id' in claims", -500, fmt.Errorf("missing or invalid 'status_id'"),
		))
	}

	if status_id == 3 {
		return c.Status(http.StatusForbidden).JSON(response.NewResponseError(
			"account_restricted", -500, fmt.Errorf("player is restricted and cannot access the website"),
		))
	} else if status_id == 4 {
		return c.Status(http.StatusForbidden).JSON(response.NewResponseError(
			"account_banned", -500, fmt.Errorf("player is banned for bad behavior and cannot access the website"),
		))
	}

	// // Validate login session
	// sv := auth.NewAuthService(connPool)
	// update_required, err := sv.CheckSession(login_session, uint64(token_version))
	// if err != nil {
	// 	smg_error := response.NewResponseError(
	// 		utls.Translate("loginSessionInvalid", nil, c),
	// 		-500,
	// 		utls.Translate("loginSessionInvalid", nil, c)),
	// 	)
	// 	return c.Status(http.StatusUnprocessableEntity).JSON(smg_error)
	// }

	// // if the token need to refresh return the update required status
	// if update_required {
	// 	smg_error := response.NewResponseError(
	// 		utls.Translate("token_update_required", nil, c),
	// 		-500,
	// 		utls.Translate("token_update_required", nil, c)),
	// 	)
	// 	return c.Status(http.StatusUpgradeRequired).JSON(smg_error)
	// }

	// Create and populate PlayerContext struct
	uCtx := types.PlayerContext{
		Id:           id,
		PlayerUuid:   playerUuid,
		UserName:     userName,
		LoginSession: login_session,
		MembershipId: membershipId,
		Exp:          time.Unix(int64(exp), 0),
		UserAgent:    string(c.Context().UserAgent()),
		Ip:           string(c.Context().RemoteIP().String()),
		StatusId:     status_id,
		TokenVersion: token_version,
	}
	c.Locals("PlayerContext", uCtx)

	return c.Next()
}
