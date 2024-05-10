package vwap

const interval = 600

// VWAPData represents the specific token's VWAP data.
type VWAPData struct {
	TokenName string  `json:"token_name"` // VWAP data belongs token name
	VWAP      float64 `json:"vwap"`       // calculated VWAP value
	Timestamp int     `json:"timestamp"`  // timestamp of the VWAP value (UNIX, 10 min interval)
}

var vwapDataMap map[string][]VWAPData

func init() {
	vwapDataMap = make(map[string][]VWAPData)
}

// store stores the VWAP data for the token or updates the existing data.
//
// Parameters:
//
//   - tokenName: the token name
//   - vwap: the VWAP value to store
//   - timestamp: the timestamp of the VWAP value
func store(tokenName string, vwap float64, timestamp int) {
	// adjust the timestamp to the 10 minutes interval.
	adjustedTimestamp := timestamp - (timestamp % interval)

	// get the VWAP data for the token
	lst, ok := vwapDataMap[tokenName]
	if !ok {
		lst = []VWAPData{}
	}

	// check last VWAP data for the list
	if len(lst) > 0 {
		last := lst[len(lst)-1]
		if last.Timestamp == adjustedTimestamp {
			last.VWAP = vwap
			vwapDataMap[tokenName] = lst
			return
		}
	}

	lst = append(lst, VWAPData{tokenName, vwap, adjustedTimestamp})
	vwapDataMap[tokenName] = lst
}
