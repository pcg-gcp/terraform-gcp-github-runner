package gcp

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func AddSecret(projectId, name, value, serviceAccount string, ctx context.Context) error {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup client: %v", err)
	}
	defer client.Close()

	replication := &secretmanagerpb.Replication{}
	if true {
		replication.Replication = &secretmanagerpb.Replication_Automatic_{
			Automatic: &secretmanagerpb.Replication_Automatic{},
		}
	} else {
		replication.Replication = &secretmanagerpb.Replication_UserManaged_{
			UserManaged: &secretmanagerpb.Replication_UserManaged{
				Replicas: []*secretmanagerpb.Replication_UserManaged_Replica{
					{
						Location: "us-central1",
					},
				},
			},
		}
	}

	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", projectId),
		SecretId: name,
		Secret: &secretmanagerpb.Secret{
			Replication: replication,
		},
	}

	secret, err := client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		return fmt.Errorf("failed to create secret: %v", err)
	}

	payload := []byte(value)
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	_, err = client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %v", err)
	}

	handle := client.IAM(secret.Name)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return fmt.Errorf("failed to get policy: %w", err)
	}

	policy.Add(fmt.Sprintf("serviceAccount:%s", serviceAccount), "roles/secretmanager.secretAccessor")
	if err = handle.SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("failed to save policy: %w", err)
	}

	return nil
}

func DeleteSecret(projectId, name string, ctx context.Context) error {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup client: %v", err)
	}
	defer client.Close()
	secretName := fmt.Sprintf("projects/%s/secrets/%s", projectId, name)
	err = client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{Name: secretName})
	if err != nil {
		return fmt.Errorf("failed to delete secret: %v", err)
	}
	return nil
}
