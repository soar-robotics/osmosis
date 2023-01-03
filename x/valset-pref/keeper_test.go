package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/osmosis-labs/osmosis/v13/app/apptesting"
	appParams "github.com/osmosis-labs/osmosis/v13/app/params"
	lockuptypes "github.com/osmosis-labs/osmosis/v13/x/lockup/types"
	"github.com/osmosis-labs/osmosis/v13/x/valset-pref/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()
}

// PrepareDelegateToValidatorSet generates 4 validators for the valsetpref.
// We self assign weights and round up to 2 decimal places in validateBasic.
func (suite *KeeperTestSuite) PrepareDelegateToValidatorSet() []types.ValidatorPreference {
	valAddrs := suite.SetupMultipleValidators(4)
	valPreferences := []types.ValidatorPreference{
		{
			ValOperAddress: valAddrs[0],
			Weight:         sdk.NewDecWithPrec(2, 1), // 0.2
		},
		{
			ValOperAddress: valAddrs[1],
			Weight:         sdk.NewDecWithPrec(332, 3), // 0.33
		},
		{
			ValOperAddress: valAddrs[2],
			Weight:         sdk.NewDecWithPrec(12, 2), // 0.12
		},
		{
			ValOperAddress: valAddrs[3],
			Weight:         sdk.NewDecWithPrec(348, 3), // 0.35
		},
	}

	return valPreferences
}

func (suite *KeeperTestSuite) GetDelegationRewards(ctx sdk.Context, valAddrStr string, delegator sdk.AccAddress) (sdk.DecCoins, stakingtypes.Validator) {
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	suite.Require().NoError(err)

	validator, found := suite.App.StakingKeeper.GetValidator(ctx, valAddr)
	suite.Require().True(found)

	endingPeriod := suite.App.DistrKeeper.IncrementValidatorPeriod(ctx, validator)

	delegation, found := suite.App.StakingKeeper.GetDelegation(ctx, delegator, valAddr)
	suite.Require().True(found)

	rewards := suite.App.DistrKeeper.CalculateDelegationRewards(ctx, validator, delegation, endingPeriod)

	return rewards, validator
}

func (suite *KeeperTestSuite) SetupExistingValidatorDelegations(ctx sdk.Context, valAddrStr string, delegator sdk.AccAddress, delegateAmt sdk.Int) {
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	suite.Require().NoError(err)

	validator, found := suite.App.StakingKeeper.GetValidator(ctx, valAddr)
	suite.Require().True(found)

	_, err = suite.App.StakingKeeper.Delegate(ctx, delegator, delegateAmt, stakingtypes.Unbonded, validator, true)
	suite.Require().NoError(err)

}

func (suite *KeeperTestSuite) SetupDelegationReward(ctx sdk.Context, delegator sdk.AccAddress, preferences []types.ValidatorPreference, existingValAddrStr string, setValSetDel, setExistingdel bool) {
	// incrementing the blockheight by 1 for reward
	ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)

	if setValSetDel {
		// only necessary if there are tokens delegated
		for _, val := range preferences {
			suite.AllocateRewards(ctx, delegator, val.ValOperAddress)
		}
	}

	if setExistingdel {
		suite.AllocateRewards(ctx, delegator, existingValAddrStr)
	}
}

func (suite *KeeperTestSuite) AllocateRewards(ctx sdk.Context, delegator sdk.AccAddress, valAddrStr string) {
	// check that there is enough reward to withdraw
	_, validator := suite.GetDelegationRewards(ctx, valAddrStr, delegator)

	// allocate some rewards
	tokens := sdk.NewDecCoins(sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10))
	suite.App.DistrKeeper.AllocateTokensToValidator(ctx, validator, tokens)

	rewardsAfterAllocation, _ := suite.GetDelegationRewards(ctx, valAddrStr, delegator)
	suite.Require().NotNil(rewardsAfterAllocation)
	suite.Require().NotZero(rewardsAfterAllocation[0].Amount)
}

func (suite *KeeperTestSuite) SetupLocks(delegator sdk.AccAddress) []lockuptypes.PeriodLock {
	locks := []lockuptypes.PeriodLock{}
	// Setup lock
	coinToLock := sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10_000_000)}
	osmoToLock := sdk.Coins{sdk.NewInt64Coin(appParams.BaseCoinUnit, 10_000_000)}
	suite.FundAcc(delegator, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000), sdk.NewInt64Coin(appParams.BaseCoinUnit, 100_000_000)})

	// happy lock case
	TwoWeekDuration, _ := time.ParseDuration("336h")
	workingLoc, err := suite.App.LockupKeeper.CreateLock(suite.Ctx, delegator, osmoToLock, TwoWeekDuration)
	suite.Require().NoError(err)

	locks = append(locks, workingLoc)

	// locking with stake denom instead of osmo denom
	stakeDenomLock, err := suite.App.LockupKeeper.CreateLock(suite.Ctx, delegator, coinToLock, TwoWeekDuration)
	suite.Require().NoError(err)

	locks = append(locks, stakeDenomLock)

	// lock case where lock owner != delegation owner
	suite.FundAcc(sdk.AccAddress([]byte("addr5---------------")), sdk.Coins{sdk.NewInt64Coin(appParams.BaseCoinUnit, 100_000_000)})
	lockWithDifferentOwner, err := suite.App.LockupKeeper.CreateLock(suite.Ctx, sdk.AccAddress([]byte("addr5---------------")), osmoToLock, TwoWeekDuration)
	suite.Require().NoError(err)

	locks = append(locks, lockWithDifferentOwner)

	// lock case where the duration != <= 2 weeks
	MorethanTwoWeekDuration, _ := time.ParseDuration("337h")
	maxDurationLock, err := suite.App.LockupKeeper.CreateLock(suite.Ctx, delegator, osmoToLock, MorethanTwoWeekDuration)
	suite.Require().NoError(err)

	locks = append(locks, maxDurationLock)

	// unbonding locks
	unbondingLocks, err := suite.App.LockupKeeper.CreateLock(suite.Ctx, delegator, osmoToLock, TwoWeekDuration)
	suite.Require().NoError(err)

	err = suite.App.LockupKeeper.BeginUnlock(suite.Ctx, unbondingLocks.ID, nil)
	suite.Require().NoError(err)

	locks = append(locks, unbondingLocks)

	// synthetic locks
	syntheticLocks, err := suite.App.LockupKeeper.CreateLock(suite.Ctx, delegator, osmoToLock, TwoWeekDuration)
	suite.Require().NoError(err)

	err = suite.App.LockupKeeper.CreateSyntheticLockup(suite.Ctx, syntheticLocks.ID, "uosmo", time.Minute, true)
	suite.Require().NoError(err)

	locks = append(locks, syntheticLocks)

	return locks
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
