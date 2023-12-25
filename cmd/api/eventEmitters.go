package main

// ClerkEventType is implemented with all the supported events described in the docs. It is overkill
// however if I need to pull from a new event I won't need to check if I've created the constant for it.
// https://clerk.com/docs/integrations/webhooks#supported-webhook-events
type ClerkEventType string

const (
	EmailCreated ClerkEventType = "email.created"

	OrganizationCreated ClerkEventType = "organization.created"
	OrganizationDeleted ClerkEventType = "organization.deleted"
	OrganizationUpdated ClerkEventType = "organization.updated"

	OrganizationInvitationAccepted ClerkEventType = "organizationInvitation.accepted"
	OrganizationInvitationCreated  ClerkEventType = "organizationInvitation.created"
	OrganizationInvitationRevoked  ClerkEventType = "organizationInvitation.revoked"

	OrganizationMembershipCreated ClerkEventType = "organizationMembership.created"
	OrganizationMembershipDeleted ClerkEventType = "organizationMembership.deleted"
	OrganizationMembershipUpdated ClerkEventType = "organizationMembership.updated"

	OrganizationDomainCreated ClerkEventType = "organizationDomain.created"
	OrganizationDomainDeleted ClerkEventType = "organizationDomain.deleted"
	OrganizationDomainUpdated ClerkEventType = "organizationDomain.updated"

	SessionCreated ClerkEventType = "session.created"
	SessionEnded   ClerkEventType = "session.ended"
	SessionRemoved ClerkEventType = "session.removed"
	SessionRevoked ClerkEventType = "session.revoked"

	SMSCreated ClerkEventType = "sms.created"

	UserCreated ClerkEventType = "user.created"
	UserDeleted ClerkEventType = "user.deleted"
	UserUpdated ClerkEventType = "user.updated"
)

type ClerkEventEmitter struct {
	handlers map[ClerkEventType][]func(data interface{})
}

func NewClerkEventEmitter() *ClerkEventEmitter {
	return &ClerkEventEmitter{
		handlers: make(map[ClerkEventType][]func(interface{})),
	}
}

func (e *ClerkEventEmitter) Register(eventType ClerkEventType, handler func(interface{})) {
	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

func (e *ClerkEventEmitter) Trigger(eventType ClerkEventType, data interface{}) {
	if handlers, ok := e.handlers[eventType]; ok {
		for _, handler := range handlers {
			handler(data)
		}
	}
}
