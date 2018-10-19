package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/ws"
)

func GetRouter(auth *AuthMiddleware) (rest.App, error) {
	return rest.MakeRouter(
		rest.Get("/api/user", auth.RequireAuth(GetUser)),
		rest.Put("/api/user", PutUser),
		rest.Post("/api/user/verify", VerifyUser),
		rest.Post("/api/user/reset-password", ResetPassword),
		rest.Post("/api/users", PostUser),
		rest.Post("/api/orgs", auth.RequireAuth(PostOrg)),
		rest.Get("/api/orgs", auth.RequireAuth(GetOrgs)),
		rest.Get("/api/orgs/:orgId", auth.RequireAuth(GetOrg)),
		rest.Put("/api/orgs/:orgId", auth.RequireAuth(PutOrg)),
		rest.Get("/api/orgs/:orgId/ledgers", auth.RequireAuth(GetOrgAccounts)),
		rest.Post("/api/orgs/:orgId/ledgers", auth.RequireAuth(PostAccount)),
		rest.Put("/api/orgs/:orgId/ledgers/:accountId", auth.RequireAuth(PutAccount)),
		rest.Delete("/api/orgs/:orgId/ledgers/:accountId", auth.RequireAuth(DeleteAccount)),
		rest.Get("/api/orgs/:orgId/ledgers/:accountId/transactions", auth.RequireAuth(GetTransactionsByAccount)),
		rest.Get("/api/orgs/:orgId/accounts", auth.RequireAuth(GetOrgAccounts)),
		rest.Post("/api/orgs/:orgId/accounts", auth.RequireAuth(PostAccount)),
		rest.Put("/api/orgs/:orgId/accounts/:accountId", auth.RequireAuth(PutAccount)),
		rest.Delete("/api/orgs/:orgId/accounts/:accountId", auth.RequireAuth(DeleteAccount)),
		rest.Get("/api/orgs/:orgId/accounts/:accountId/transactions", auth.RequireAuth(GetTransactionsByAccount)),
		rest.Get("/api/orgs/:orgId/transactions", auth.RequireAuth(GetTransactionsByOrg)),
		rest.Post("/api/orgs/:orgId/transactions", auth.RequireAuth(PostTransaction)),
		rest.Put("/api/orgs/:orgId/transactions/:transactionId", auth.RequireAuth(PutTransaction)),
		rest.Delete("/api/orgs/:orgId/transactions/:transactionId", auth.RequireAuth(DeleteTransaction)),
		rest.Get("/api/orgs/:orgId/prices", auth.RequireAuth(GetPrices)),
		rest.Post("/api/orgs/:orgId/prices", auth.RequireAuth(PostPrice)),
		rest.Delete("/api/orgs/:orgId/prices/:priceId", auth.RequireAuth(DeletePrice)),
		rest.Get("/ws", ws.Handler),
		rest.Post("/api/sessions", auth.RequireAuth(PostSession)),
		rest.Delete("/api/sessions/:sessionId", auth.RequireAuth(DeleteSession)),
		rest.Get("/api/apikeys", auth.RequireAuth(GetApiKeys)),
		rest.Post("/api/apikeys", auth.RequireAuth(PostApiKey)),
		rest.Put("/api/apikeys/:apiKeyId", auth.RequireAuth(PutApiKey)),
		rest.Delete("/api/apikeys/:apiKeyId", auth.RequireAuth(DeleteApiKey)),
		rest.Get("/api/orgs/:orgId/invites", auth.RequireAuth(GetInvites)),
		rest.Post("/api/orgs/:orgId/invites", auth.RequireAuth(PostInvite)),
		rest.Put("/api/orgs/:orgId/invites/:inviteId", auth.RequireAuth(PutInvite)),
		rest.Delete("/api/orgs/:orgId/invites/:inviteId", auth.RequireAuth(DeleteInvite)),
	)
}
