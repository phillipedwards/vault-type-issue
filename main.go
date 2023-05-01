package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-vault/sdk/v5/go/vault"
	"github.com/pulumi/pulumi-vault/sdk/v5/go/vault/approle"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		token := pulumi.String("hvs.JNtqfYVdz2Lnrbfyskcob1rE")
		address := pulumi.String("http://127.0.0.1:8200")

		provider, err := vault.NewProvider(ctx, "provider", &vault.ProviderArgs{
			Token:   token,
			Address: address,
		})

		if err != nil {
			return err
		}

		provider2, err := vault.NewProvider(ctx, "provider2", &vault.ProviderArgs{
			Token:   token,
			Address: address,
		}, pulumi.DependsOn([]pulumi.Resource{provider}))

		if err != nil {
			return err
		}

		role, err := vault.NewAuthBackend(ctx, "approle", &vault.AuthBackendArgs{
			Type: pulumi.String("approle"),
		}, pulumi.Provider(provider2))

		if err != nil {
			return err
		}

		example, err := approle.NewAuthBackendRole(ctx, "example", &approle.AuthBackendRoleArgs{
			Backend:  role.Path,
			RoleName: pulumi.String("test-role"),
			TokenPolicies: pulumi.StringArray{
				pulumi.String("default"),
				pulumi.String("dev"),
				pulumi.String("prod"),
			},
		}, pulumi.Provider(provider2))

		if err != nil {
			return err
		}

		tmpJSON0, err := json.Marshal(map[string]interface{}{
			"hello": "world",
		})

		if err != nil {
			return err
		}

		json0 := string(tmpJSON0)
		_, err = approle.NewAuthBackendRoleSecretId(ctx, "id", &approle.AuthBackendRoleSecretIdArgs{
			Backend:  role.Path,
			RoleName: example.RoleName,
			Metadata: pulumi.String(json0),
		}, pulumi.Provider(provider2))

		if err != nil {
			return err
		}

		return nil
	})
}
