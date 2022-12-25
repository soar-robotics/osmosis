package math

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TicksToSqrtPrice returns the sqrt price for the lower and upper ticks.
// Returns error if fails to calculate sqrt price.
// TODO: spec and tests
func TicksToSqrtPrice(lowerTick, upperTick int64) (sdk.Dec, sdk.Dec, error) {
	// TODO: Dont hardcode
	sqrtPriceUpperTick, err := TickToSqrtPrice(sdk.NewInt(upperTick), sdk.NewInt(-4))
	if err != nil {
		return sdk.Dec{}, sdk.Dec{}, err
	}
	sqrtPriceLowerTick, err := TickToSqrtPrice(sdk.NewInt(lowerTick), sdk.NewInt(-4))
	if err != nil {
		return sdk.Dec{}, sdk.Dec{}, err
	}
	return sqrtPriceLowerTick, sqrtPriceUpperTick, nil
}

// // TickToSqrtPrice takes the tick index and returns the corresponding sqrt of the price.
// // Returns error if fails to calculate sqrt price. Otherwise, the computed value and nil.
// // TODO: test
// func TickToSqrtPrice(tickIndex sdk.Int) (sqrtPrice sdk.Dec, err error) {
// 	if tickIndex.IsZero() {
// 		return sdk.OneDec(), nil
// 	}
// 	kAtPrice1 := sdk.NewInt(-6)

// 	var kIncrementDistance sdk.Int
// 	if kAtPrice1.GTE(sdk.ZeroInt()) {
// 		kIncrementDistance = sdk.NewDec(9).Mul(sdk.OneDec().Quo(sdk.NewInt(10).ToDec().Power(kAtPrice1.Uint64()))).TruncateInt()
// 		fmt.Printf("kIncrementDistance pos: %v \n", kIncrementDistance)
// 	} else {
// 		kIncrementDistance = sdk.NewDec(9).Mul(sdk.NewInt(10).ToDec().Power(kAtPrice1.Abs().Uint64())).TruncateInt()
// 		fmt.Printf("kIncrementDistance neg: %v \n", kIncrementDistance)
// 	}
// 	kDelta := tickIndex.ToDec().Quo(kIncrementDistance.ToDec()).TruncateInt()
// 	fmt.Printf("kDelta: %v \n", kDelta)
// 	curK := kAtPrice1.Add(kDelta)
// 	fmt.Printf("curK: %v \n", curK)

// 	if tickIndex.LT(sdk.NewInt(999)) {
// 		kIncrementDistance = sdk.NewInt(999)
// 		kDelta = sdk.NewInt(6)
// 		fmt.Printf("kDelta lt 999: %v \n", kDelta)
// 		curK = kAtPrice1.Add(kDelta)
// 		fmt.Printf("curK lt 999: %v \n", curK)
// 		kDelta = sdk.NewInt(0)
// 	}

// 	if tickIndex.GT(sdk.NewInt(999)) && tickIndex.LT(sdk.NewInt(90999)) {
// 		kIncrementDistance = sdk.NewInt(999)
// 		kDelta = sdk.NewInt(5)
// 		curK = kAtPrice1.Add(kDelta)
// 		kDelta = sdk.NewInt(1)

// 	}

// 	var curIncrement sdk.Dec
// 	if curK.GTE(sdk.ZeroInt()) {
// 		curIncrement = sdk.NewDec(10).Power(curK.Uint64())
// 		fmt.Printf("curIncrement ps: %v \n", curIncrement)
// 	} else {
// 		curIncrement = sdk.NewDec(1).Quo(sdk.NewDec(10).Power(curK.Abs().Uint64()))
// 		fmt.Printf("curIncrement ng: %v \n", curIncrement)
// 	}

// 	pre := kDelta.Mul(kIncrementDistance)
// 	fmt.Printf("pre: %v \n", pre)
// 	numAdditiveTicks := tickIndex.Sub(pre)
// 	fmt.Printf("numAdditiveTicks: %v \n", numAdditiveTicks)
// 	var price sdk.Dec
// 	if kDelta.GTE(sdk.ZeroInt()) {
// 		first := sdk.NewDec(10).Power(kDelta.Uint64())
// 		price = first.Add(numAdditiveTicks.ToDec().Mul(curIncrement))
// 		fmt.Printf("price pos: %v \n", price)
// 	} else {
// 		first := sdk.OneDec().Quo(sdk.NewDec(10).Power(kDelta.Abs().Uint64()))
// 		price = first.Add(numAdditiveTicks.ToDec().Mul(curIncrement))
// 		fmt.Printf("price neg: %v \n", price)
// 	}
// 	return price, nil
// }

// TickToSqrtPrice calculates the price at a given tick index based on the provided
// starting price of 1 at k=0. The price is calculated using a square root function
// with a coefficient of 9, where the price increases by a factor of 10 for every
// increment of k.
//
// The function takes in two arguments:
// 	- tickIndex: the tick index to calculate the price for
// 	- kAtPriceOne: the value of k at which the starting price of 1 is set
//
// It returns a sdk.Dec representing the calculated price and an error if any errors
// occurred during the calculation.
func TickToSqrtPrice(tickIndex, kAtPriceOne sdk.Int) (price sdk.Dec, err error) {
	if tickIndex.IsZero() {
		return sdk.OneDec(), nil
	}

	fmt.Printf("tickIndex: %v \n", tickIndex)
	fmt.Printf("kAtPriceOne: %v \n", kAtPriceOne)

	var kIncrementDistance sdk.Dec
	// The formula is as follows: k_increment_distance = 9 * 10**(-k_at_price_1)
	// Due to sdk.Power restrictions, if the resulting power is negative, we take 9 * (1/10**k_at_price_1)
	if kAtPriceOne.GTE(sdk.ZeroInt()) {
		kIncrementDistance = sdk.NewDec(9).Mul(sdk.OneDec().Quo(sdk.NewInt(10).ToDec().Power(kAtPriceOne.Uint64())))
		fmt.Printf("kIncrementDistance pos: %v \n", kIncrementDistance)
	} else {
		kIncrementDistance = sdk.NewDec(9).Mul(sdk.NewInt(10).ToDec().Power(kAtPriceOne.Abs().Uint64()))
		fmt.Printf("kIncrementDistance neg: %v \n", kIncrementDistance)
	}

	// Use floor division to determine how many k increments we have passed
	kDelta := tickIndex.ToDec().Quo(kIncrementDistance).TruncateInt()
	fmt.Printf("kDelta: %v \n", kDelta)

	// Calculate the current k value from the starting k value and the k delta
	curK := (kAtPriceOne.Add(kDelta))
	fmt.Printf("curK: %v \n", curK)

	var curIncrement sdk.Dec
	if curK.GTE(sdk.ZeroInt()) {
		curIncrement = sdk.NewDec(10).Power(curK.Uint64())
		fmt.Printf("curIncrement ps: %v \n", curIncrement)
	} else {
		curIncrement = sdk.NewDec(1).Quo(sdk.NewDec(10).Power(curK.Abs().Uint64()))
		fmt.Printf("curIncrement ng: %v \n", curIncrement)
	}

	numAdditiveTicks := tickIndex.ToDec().Sub(kDelta.ToDec().Mul(kIncrementDistance))
	fmt.Printf("numAdditiveTicks: %v \n", numAdditiveTicks)

	if kDelta.GTE(sdk.ZeroInt()) {
		price = sdk.NewDec(10).Power(kDelta.Uint64()).Add(numAdditiveTicks.Mul(curIncrement))
		fmt.Printf("price pos: %v \n", price)
	} else {
		price = sdk.OneDec().Quo(sdk.NewDec(10).Power(kDelta.Abs().Uint64())).Add(numAdditiveTicks.Mul(curIncrement))
		fmt.Printf("price neg: %v \n", price)
	}
	fmt.Println()
	return price, nil
}

// PriceToTick takes a price and returns the corresponding tick index
func PriceToTick(price sdk.Dec, kAtPriceOne sdk.Int) (tickIndex sdk.Int) {
	if price.Equal(sdk.OneDec()) {
		return sdk.ZeroInt()
	}

	var kIncrementDistance sdk.Dec
	// The formula is as follows: k_increment_distance = 9 * 10**(-k_at_price_1)
	// Due to sdk.Power restrictions, if the resulting power is negative, we take 9 * (1/10**k_at_price_1)
	if kAtPriceOne.GTE(sdk.ZeroInt()) {
		kIncrementDistance = sdk.NewDec(9).Mul(sdk.OneDec().Quo(sdk.NewInt(10).ToDec().Power(kAtPriceOne.Uint64())))
		fmt.Printf("kIncrementDistance pos: %v \n", kIncrementDistance)
	} else {
		kIncrementDistance = sdk.NewDec(9).Mul(sdk.NewInt(10).ToDec().Power(kAtPriceOne.Abs().Uint64()))
		fmt.Printf("kIncrementDistance neg: %v \n", kIncrementDistance)
	}

	total := sdk.OneDec()
	ticksPassed := sdk.ZeroInt()
	currentK := kAtPriceOne

	var curIncrement sdk.Dec
	if currentK.GTE(sdk.ZeroInt()) {
		curIncrement = sdk.NewDec(10).Power(currentK.Uint64())
	} else {
		curIncrement = sdk.NewDec(1).Quo(sdk.NewDec(10).Power(currentK.Abs().Uint64()))
	}

	fmt.Printf("total: %v \n", total)
	fmt.Printf("price: %v \n", price)

	for total.LT(price) {
		if currentK.GTE(sdk.ZeroInt()) {
			curIncrement = sdk.NewDec(10).Power(currentK.Uint64())
			fmt.Printf("curIncrement ps: %v \n", curIncrement)
			maxPriceForCurrentIncrement := kIncrementDistance.Mul(curIncrement)
			if total.Add(maxPriceForCurrentIncrement).LT(price) {
				total = total.Add(maxPriceForCurrentIncrement)
				currentK = currentK.Add(sdk.OneInt())
				ticksPassed = ticksPassed.Add(kIncrementDistance.TruncateInt())
				fmt.Printf("total ps: %v \n", total)
				fmt.Printf("newcurrentK ps: %v \n", currentK)
			} else {
				break
			}
		} else {
			curIncrement = sdk.NewDec(1).Quo(sdk.NewDec(10).Power(currentK.Abs().Uint64()))
			fmt.Printf("curIncrement ng: %v \n", curIncrement)
			maxPriceForCurrentIncrement := kIncrementDistance.Mul(curIncrement)
			if total.Add(maxPriceForCurrentIncrement).LT(price) {
				total = total.Add(maxPriceForCurrentIncrement)
				currentK = currentK.Add(sdk.OneInt())
				ticksPassed = ticksPassed.Add(kIncrementDistance.TruncateInt())
				fmt.Printf("total ng: %v \n", total)
				fmt.Printf("newcurrentK ng: %v \n", currentK)
			} else {
				break
			}
		}
	}
	ticksToBeFulfilledByCurrentK := price.Sub(total).Quo(curIncrement)

	ticksPassed = ticksPassed.Add(ticksToBeFulfilledByCurrentK.TruncateInt())

	return ticksPassed
}

// // PriceToTick takes a price and returns the corresponding tick index
// func PriceToTick(price sdk.Dec, kAtPriceOne sdk.Int) (tickIndex sdk.Int) {
// 	if price.Equal(sdk.OneDec()) {
// 		return sdk.ZeroInt()
// 	}

// 	var kIncrementDistance sdk.Dec
// 	// The formula is as follows: k_increment_distance = 9 * 10**(-k_at_price_1)
// 	// Due to sdk.Power restrictions, if the resulting power is negative, we take 9 * (1/10**k_at_price_1)
// 	if kAtPriceOne.GTE(sdk.ZeroInt()) {
// 		kIncrementDistance = sdk.NewDec(9).Mul(sdk.OneDec().Quo(sdk.NewInt(10).ToDec().Power(kAtPriceOne.Uint64())))
// 		fmt.Printf("kIncrementDistance pos: %v \n", kIncrementDistance)
// 	} else {
// 		kIncrementDistance = sdk.NewDec(9).Mul(sdk.NewInt(10).ToDec().Power(kAtPriceOne.Abs().Uint64()))
// 		fmt.Printf("kIncrementDistance neg: %v \n", kIncrementDistance)
// 	}

// 	total := sdk.OneDec()
// 	ticksPassed := sdk.ZeroInt()
// 	currentK := kAtPriceOne
// 	var curIncrement sdk.Dec

// 	fmt.Printf("total: %v \n", total)
// 	fmt.Printf("price: %v \n", price)

// 	for total.LT(price) {
// 		if currentK.GTE(sdk.ZeroInt()) {
// 			curIncrement = sdk.NewDec(10).Power(currentK.Uint64())
// 			fmt.Printf("curIncrement ps: %v \n", curIncrement)
// 			maxPriceForCurrentIncrement := kIncrementDistance.Mul(curIncrement)
// 			if total.Add(maxPriceForCurrentIncrement).LT(price) {
// 				total = total.Add(maxPriceForCurrentIncrement)
// 				currentK = currentK.Add(sdk.OneInt())
// 				ticksPassed = ticksPassed.Add(kIncrementDistance.TruncateInt())
// 				fmt.Printf("total ps: %v \n", total)
// 				fmt.Printf("newcurrentK ps: %v \n", currentK)
// 			} else {
// 				break
// 			}
// 		} else {
// 			curIncrement = sdk.NewDec(1).Quo(sdk.NewDec(10).Power(currentK.Abs().Uint64()))
// 			fmt.Printf("curIncrement ng: %v \n", curIncrement)
// 			maxPriceForCurrentIncrement := kIncrementDistance.Mul(curIncrement)
// 			if total.Add(maxPriceForCurrentIncrement).LT(price) {
// 				total = total.Add(maxPriceForCurrentIncrement)
// 				currentK = currentK.Add(sdk.OneInt())
// 				ticksPassed = ticksPassed.Add(kIncrementDistance.TruncateInt())
// 				fmt.Printf("total ng: %v \n", total)
// 				fmt.Printf("newcurrentK ng: %v \n", currentK)
// 			} else {
// 				break
// 			}
// 		}
// 	}
// 	ticksToBeFulfilledByCurrentK := price.Sub(total).Quo(curIncrement)

// 	ticksPassed = ticksPassed.Add(ticksToBeFulfilledByCurrentK.TruncateInt())

// 	return ticksPassed
// }
