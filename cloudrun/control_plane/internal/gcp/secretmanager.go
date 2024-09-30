package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func (c *Client) addSecret(name, value, serviceAccount string, ctx context.Context) error {
	replication := &secretmanagerpb.Replication{}
	if true {
		replication.Replication = &secretmanagerpb.Replication_UserManaged_{
			UserManaged: &secretmanagerpb.Replication_UserManaged{
				Replicas: []*secretmanagerpb.Replication_UserManaged_Replica{
					{
						Location: c.cfg.Region,
					},
				},
			},
		}
	} else {
		replication.Replication = &secretmanagerpb.Replication_Automatic_{
			Automatic: &secretmanagerpb.Replication_Automatic{},
		}
	}

	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", c.cfg.ProjectID),
		SecretId: name,
		Secret: &secretmanagerpb.Secret{
			Replication: replication,
		},
	}

	secret, err := c.secretManagerClient.CreateSecret(ctx, createSecretReq)
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

	_, err = c.secretManagerClient.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %v", err)
	}

	handle := c.secretManagerClient.IAM(secret.Name)
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

func (c *Client) deleteSecret(name string, ctx context.Context) error {
	secretName := fmt.Sprintf("projects/%s/secrets/%s", c.cfg.ProjectID, name)
	err := c.secretManagerClient.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{Name: secretName})
	if err != nil {
		return fmt.Errorf("failed to delete secret: %v", err)
	}
	return nil
}
