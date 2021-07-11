// Code generated by protoc-gen-entgrpc. DO NOT EDIT.
package entpb

import (
	context "context"
	ent "entgo.io/contrib/entproto/internal/todo/ent"
	attachment "entgo.io/contrib/entproto/internal/todo/ent/attachment"
	user "entgo.io/contrib/entproto/internal/todo/ent/user"
	sqlgraph "entgo.io/ent/dialect/sql/sqlgraph"
	uuid "github.com/google/uuid"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// AttachmentService implements AttachmentServiceServer
type AttachmentService struct {
	client *ent.Client
	UnimplementedAttachmentServiceServer
}

// NewAttachmentService returns a new AttachmentService
func NewAttachmentService(client *ent.Client) *AttachmentService {
	return &AttachmentService{
		client: client,
	}
}

// toProtoAttachment transforms the ent type to the pb type
func toProtoAttachment(e *ent.Attachment) (*Attachment, error) {
	v := &Attachment{}
	id, err := e.ID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	v.Id = id

	for _, edg := range e.Edges.Recipients {

		id := int32(edg.ID)

		v.Recipients = append(v.Recipients, &User{
			Id: id,
		})
	}

	if edg := e.Edges.User; edg != nil {

		id := int32(edg.ID)

		v.User = &User{
			Id: id,
		}
	}

	return v, nil
}

// Create implements AttachmentServiceServer.Create
func (svc *AttachmentService) Create(ctx context.Context, req *CreateAttachmentRequest) (*Attachment, error) {
	attachment := req.GetAttachment()
	m := svc.client.Attachment.Create()
	for _, item := range attachment.GetRecipients() {
		recipients := int(item.GetId())
		m.AddRecipientIDs(recipients)
	}

	attachmentUser := int(attachment.GetUser().GetId())
	m.SetUserID(attachmentUser)
	res, err := m.Save(ctx)

	switch {
	case err == nil:
		proto, err := toProtoAttachment(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal: %s", err)
	}
}

// Get implements AttachmentServiceServer.Get
func (svc *AttachmentService) Get(ctx context.Context, req *GetAttachmentRequest) (*Attachment, error) {
	var (
		err error
		get *ent.Attachment
	)
	var id uuid.UUID
	if err := (&id).UnmarshalBinary(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}

	switch req.GetView() {
	case GetAttachmentRequest_VIEW_UNSPECIFIED, GetAttachmentRequest_BASIC:
		get, err = svc.client.Attachment.Get(ctx, id)
	case GetAttachmentRequest_WITH_EDGE_IDS:
		get, err = svc.client.Attachment.Query().
			Where(attachment.ID(id)).
			WithRecipients(func(query *ent.UserQuery) {
				query.Select(user.FieldID)
			}).
			WithUser(func(query *ent.UserQuery) {
				query.Select(user.FieldID)
			}).
			Only(ctx)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProtoAttachment(get)
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}
}

// Update implements AttachmentServiceServer.Update
func (svc *AttachmentService) Update(ctx context.Context, req *UpdateAttachmentRequest) (*Attachment, error) {
	attachment := req.GetAttachment()
	var attachmentID uuid.UUID
	if err := (&attachmentID).UnmarshalBinary(attachment.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}

	m := svc.client.Attachment.UpdateOneID(attachmentID)
	for _, item := range attachment.GetRecipients() {
		recipients := int(item.GetId())
		m.AddRecipientIDs(recipients)
	}

	attachmentUser := int(attachment.GetUser().GetId())
	m.SetUserID(attachmentUser)
	res, err := m.Save(ctx)

	switch {
	case err == nil:
		proto, err := toProtoAttachment(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal: %s", err)
	}
}

// Delete implements AttachmentServiceServer.Delete
func (svc *AttachmentService) Delete(ctx context.Context, req *DeleteAttachmentRequest) (*emptypb.Empty, error) {
	var err error
	var id uuid.UUID
	if err := (&id).UnmarshalBinary(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}

	err = svc.client.Attachment.DeleteOneID(id).Exec(ctx)
	switch {
	case err == nil:
		return &emptypb.Empty{}, nil
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}
}
