// SPDX-License-Identifier: Apache-2.0

package updatepolicythreshold

import (
	"os"

	"github.com/gittuf/gittuf/internal/cmd/common"
	"github.com/gittuf/gittuf/internal/cmd/trust/persistent"
	"github.com/gittuf/gittuf/internal/repository"
	"github.com/spf13/cobra"
)

type options struct {
	p         *persistent.Options
	threshold int
}

func (o *options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(
		&o.threshold,
		"threshold",
		-1,
		"threshold of valid signatures required for main policy",
	)
	cmd.MarkFlagRequired("threshold") //nolint:errcheck
}

func (o *options) Run(cmd *cobra.Command, _ []string) error {
	repo, err := repository.LoadRepository()
	if err != nil {
		return err
	}

	rootKeyBytes, err := os.ReadFile(o.p.SigningKey)
	if err != nil {
		return err
	}

	return repo.UpdateTopLevelTargetsThreshold(cmd.Context(), rootKeyBytes, o.threshold, true)
}

func New(persistent *persistent.Options) *cobra.Command {
	o := &options{p: persistent}
	cmd := &cobra.Command{
		Use:     "update-policy-threshold",
		Short:   "Update Policy threshold in the gittuf root of trust",
		Long:    `This command allows users to update the threshold of valid signatures required for the policy. DO NOT USE until policy-staging is working, so that multiple developers can sequentially sign the policy metadata.`,
		PreRunE: common.CheckIfSigningViable,
		RunE:    o.Run,
	}
	o.AddFlags(cmd)

	return cmd
}
