package twapcli

import (
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/osmosis-labs/osmosis/osmoutils/osmocli"
	"github.com/osmosis-labs/osmosis/v13/x/twap/client/queryproto"
	"github.com/osmosis-labs/osmosis/v13/x/twap/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := osmocli.QueryIndexCmd(types.ModuleName)
	cmd.AddCommand(GetQueryTwapCommand())

	return cmd
}

// GetQueryTwapCommand returns multiplier of an asset by denom.
func GetQueryTwapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "twap [poolid] [base denom] [start time] [end time]",
		Short: "Query twap",
		Long: osmocli.FormatLongDescDirect(`Query twap for pool. Start time must be unix time. End time can be unix time or duration.

Example:
{{.CommandPrefix}} twap 1 uosmo 1667088000 24h
{{.CommandPrefix}} twap 1 uosmo 1667088000 1667174400
`, types.ModuleName),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// boilerplate parse fields
			poolId, baseDenom, startTime, endTime, err := twapQueryParseArgs(args)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			quoteDenom, err := getQuoteDenomFromLiquidity(cmd.Context(), clientCtx, poolId, baseDenom)
			if err != nil {
				return err
			}

			queryClient := queryproto.NewQueryClient(clientCtx)
			res, err := queryClient.ArithmeticTwap(cmd.Context(), &queryproto.ArithmeticTwapRequest{
				PoolId:     poolId,
				BaseAsset:  baseDenom,
				QuoteAsset: quoteDenom,
				StartTime:  startTime,
				EndTime:    &endTime,
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func twapQueryParseArgs(args []string) (poolId uint64, baseDenom string, startTime time.Time, endTime time.Time, err error) {
	// boilerplate parse fields
	// <UINT PARSE>
	poolId, err = osmocli.ParseUint(args[0], "poolId")
	if err != nil {
		return
	}

	// <DENOM PARSE>
	baseDenom = strings.TrimSpace(args[1])

	// <UNIX TIME PARSE>
	startTime, err = osmocli.ParseUnixTime(args[2], "start time")
	if err != nil {
		return
	}

	// END TIME PARSE: ONEOF {<UNIX TIME PARSE>, <DURATION>}
	// try parsing in unix time, if failed try parsing in duration
	endTime, err = osmocli.ParseUnixTime(args[3], "end time")
	if err != nil {
		// TODO if we don't use protoreflect:
		// make better error combiner, rather than just returning last error
		duration, err2 := time.ParseDuration(args[3])
		if err2 != nil {
			err = err2
			return
		}
		endTime = startTime.Add(duration)
	}
	return poolId, baseDenom, startTime, endTime, nil
}
