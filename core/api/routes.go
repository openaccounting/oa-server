package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/ws"
)

func GetRouter(auth *AuthMiddleware, prefix string) (rest.App, error) {
	return rest.MakeRouter(
		rest.Get(prefix+"/user", auth.RequireAuth(GetUser)),
		rest.Put(prefix+"/user", PutUser),
		rest.Post(prefix+"/user/verify", VerifyUser),
		rest.Post(prefix+"/user/reset-password", ResetPassword),
		rest.Post(prefix+"/users", PostUser),
		rest.Post(prefix+"/orgs", auth.RequireAuth(PostOrg)),
		rest.Get(prefix+"/orgs", auth.RequireAuth(GetOrgs)),
		rest.Get(prefix+"/orgs/:orgId", auth.RequireAuth(GetOrg)),
		rest.Put(prefix+"/orgs/:orgId", auth.RequireAuth(PutOrg)),
		rest.Get(prefix+"/orgs/:orgId/ledgers", auth.RequireAuth(GetOrgAccounts)),
		rest.Post(prefix+"/orgs/:orgId/ledgers", auth.RequireAuth(PostAccount)),
		rest.Put(prefix+"/orgs/:orgId/ledgers/:accountId", auth.RequireAuth(PutAccount)),
		rest.Delete(prefix+"/orgs/:orgId/ledgers/:accountId", auth.RequireAuth(DeleteAccount)),
		rest.Get(prefix+"/orgs/:orgId/ledgers/:accountId/transactions", auth.RequireAuth(GetTransactionsByAccount)),
		rest.Get(prefix+"/orgs/:orgId/accounts", auth.RequireAuth(GetOrgAccounts)),
		rest.Post(prefix+"/orgs/:orgId/accounts", auth.RequireAuth(PostAccount)),
		rest.Put(prefix+"/orgs/:orgId/accounts/:accountId", auth.RequireAuth(PutAccount)),
		rest.Delete(prefix+"/orgs/:orgId/accounts/:accountId", auth.RequireAuth(DeleteAccount)),
		rest.Get(prefix+"/orgs/:orgId/accounts/:accountId/transactions", auth.RequireAuth(GetTransactionsByAccount)),
		rest.Get(prefix+"/orgs/:orgId/transactions", auth.RequireAuth(GetTransactionsByOrg)),
		rest.Post(prefix+"/orgs/:orgId/transactions", auth.RequireAuth(PostTransaction)),
		rest.Put(prefix+"/orgs/:orgId/transactions/:transactionId", auth.RequireAuth(PutTransaction)),
		rest.Delete(prefix+"/orgs/:orgId/transactions/:transactionId", auth.RequireAuth(DeleteTransaction)),
		rest.Get(prefix+"/orgs/:orgId/prices", auth.RequireAuth(GetPrices)),
		rest.Post(prefix+"/orgs/:orgId/prices", auth.RequireAuth(PostPrice)),
		rest.Delete(prefix+"/orgs/:orgId/prices/:priceId", auth.RequireAuth(DeletePrice)),
		rest.Get("/ws", ws.Handler),
		rest.Post(prefix+"/sessions", auth.RequireAuth(PostSession)),
		rest.Delete(prefix+"/sessions/:sessionId", auth.RequireAuth(DeleteSession)),
		rest.Get(prefix+"/apikeys", auth.RequireAuth(GetApiKeys)),
		rest.Post(prefix+"/apikeys", auth.RequireAuth(PostApiKey)),
		rest.Put(prefix+"/apikeys/:apiKeyId", auth.RequireAuth(PutApiKey)),
		rest.Delete(prefix+"/apikeys/:apiKeyId", auth.RequireAuth(DeleteApiKey)),
		rest.Get(prefix+"/orgs/:orgId/invites", auth.RequireAuth(GetInvites)),
		rest.Post(prefix+"/orgs/:orgId/invites", auth.RequireAuth(PostInvite)),
		rest.Put(prefix+"/orgs/:orgId/invites/:inviteId", auth.RequireAuth(PutInvite)),
		rest.Delete(prefix+"/orgs/:orgId/invites/:inviteId", auth.RequireAuth(DeleteInvite)),
		rest.Get(prefix+"/health-check", GetSystemHealthStatus),
	)
}
