package transaction

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"net/url"
)

type DeepLink struct {
	Type    Type
	Data    []byte
	Payload []byte

	Nonce    *uint // optional
	GasPrice *uint // optional
	GasCoin  *Coin // optional
}

func (d *DeepLink) CreateLink(pass string) (string, error) {
	tx, err := d.Encode()
	if err != nil {
		return "", err
	}

	rawQuery := ""
	if pass != "" {
		rawQuery = fmt.Sprintf("p=%s", base64.RawStdEncoding.EncodeToString([]byte(pass)))
	}

	u := &url.URL{
		Scheme:   "https",
		Host:     "bip.to",
		Path:     fmt.Sprintf("/tx/%s", tx),
		RawQuery: rawQuery,
	}
	return u.String(), nil
}

func (d *DeepLink) Encode() (string, error) {
	src, err := rlp.EncodeToBytes(d)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(src), nil
}

func (d *DeepLink) setType(t Type) *DeepLink {
	d.Type = t
	return d
}

func (d *DeepLink) SetPayload(payload []byte) *DeepLink {
	d.Payload = payload
	return d
}

func (d *DeepLink) SetGasCoin(symbol string) *DeepLink {
	gasCoin := Coin{}
	d.GasCoin = &gasCoin
	copy(d.GasCoin[:], symbol)
	return d
}

func NewDeepLink(data DataInterface) (*DeepLink, error) {
	d := new(DeepLink)

	bytes, err := data.encode()
	if err != nil {
		return d, err
	}
	d.Data = bytes

	switch data.(type) {
	case *SendData:
		return d.setType(TypeSend), nil
	case *SellCoinData:
		return d.setType(TypeSellCoin), nil
	case *SellAllCoinData:
		return d.setType(TypeSellAllCoin), nil
	case *BuyCoinData:
		return d.setType(TypeBuyCoin), nil
	case *CreateCoinData:
		return d.setType(TypeCreateCoin), nil
	case *DeclareCandidacyData:
		return d.setType(TypeDeclareCandidacy), nil
	case *DelegateData:
		return d.setType(TypeDelegate), nil
	case *UnbondData:
		return d.setType(TypeUnbond), nil
	case *RedeemCheckData:
		return d.setType(TypeRedeemCheck), nil
	case *SetCandidateOnData:
		return d.setType(TypeSetCandidateOnline), nil
	case *SetCandidateOffData:
		return d.setType(TypeSetCandidateOffline), nil
	case *CreateMultisigData:
		return d.setType(TypeCreateMultisig), nil
	case *MultisendData:
		return d.setType(TypeMultisend), nil
	case *EditCandidateData:
		return d.setType(TypeEditCandidate), nil

	default:
		return nil, errors.New("unknown transaction type")
	}
}
