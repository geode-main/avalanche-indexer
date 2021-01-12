package stream

import (
	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/codec/hierarchycodec"
	"github.com/ava-labs/avalanchego/codec/linearcodec"
	"github.com/ava-labs/avalanchego/utils/codec"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/avalanchego/vms/avm"
	"github.com/ava-labs/avalanchego/vms/nftfx"
	"github.com/ava-labs/avalanchego/vms/propertyfx"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
)

// InitCodecManager returns a configured codec manager with support for multiple codecs
func InitCodecManager() (codec.Manager, error) {
	errs := wrappers.Errs{}
	preApricotCodecVersion := uint16(0)
	apricotCodecVersion := uint16(1)

	preApricotCodec := initPreApricotCodec(&errs)
	apricotCodec := initApricotCodec(&errs)

	if errs.Errored() {
		return nil, errs.Err()
	}

	manager := codec.NewManager()
	manager.RegisterCodec(preApricotCodecVersion, preApricotCodec)
	manager.RegisterCodec(apricotCodecVersion, apricotCodec)

	return manager, nil
}

func initPreApricotCodec(errs *wrappers.Errs) linearcodec.Codec {
	c := linearcodec.NewDefault()

	errs.Add(
		c.RegisterType(&avm.BaseTx{}),
		c.RegisterType(&avm.CreateAssetTx{}),
		c.RegisterType(&avm.OperationTx{}),
		c.RegisterType(&avm.ImportTx{}),
		c.RegisterType(&avm.ExportTx{}),
		c.RegisterType(&secp256k1fx.TransferInput{}),
		c.RegisterType(&secp256k1fx.MintOutput{}),
		c.RegisterType(&secp256k1fx.TransferOutput{}),
		c.RegisterType(&secp256k1fx.MintOperation{}),
		c.RegisterType(&secp256k1fx.Credential{}),
		c.RegisterType(&nftfx.MintOutput{}),
		c.RegisterType(&nftfx.TransferOutput{}),
		c.RegisterType(&nftfx.MintOperation{}),
		c.RegisterType(&nftfx.TransferOperation{}),
		c.RegisterType(&nftfx.Credential{}),
		c.RegisterType(&propertyfx.MintOutput{}),
		c.RegisterType(&propertyfx.OwnedOutput{}),
		c.RegisterType(&propertyfx.MintOperation{}),
		c.RegisterType(&propertyfx.BurnOperation{}),
		c.RegisterType(&propertyfx.Credential{}),
	)

	return c
}

func initApricotCodec(errs *wrappers.Errs) hierarchycodec.Codec {
	c := hierarchycodec.NewDefault()

	errs.Add(
		c.RegisterType(&avm.BaseTx{}),
		c.RegisterType(&avm.CreateAssetTx{}),
		c.RegisterType(&avm.OperationTx{}),
		c.RegisterType(&avm.ImportTx{}),
		c.RegisterType(&avm.ExportTx{}),
	)
	c.NextGroup()

	errs.Add(
		c.RegisterType(&secp256k1fx.TransferInput{}),
		c.RegisterType(&secp256k1fx.MintOutput{}),
		c.RegisterType(&secp256k1fx.TransferOutput{}),
		c.RegisterType(&secp256k1fx.MintOperation{}),
		c.RegisterType(&secp256k1fx.Credential{}),
		c.RegisterType(&secp256k1fx.ManagedAssetStatusOutput{}),
		c.RegisterType(&secp256k1fx.UpdateManagedAssetOperation{}),
	)
	c.NextGroup()

	errs.Add(
		c.RegisterType(&nftfx.MintOutput{}),
		c.RegisterType(&nftfx.TransferOutput{}),
		c.RegisterType(&nftfx.MintOperation{}),
		c.RegisterType(&nftfx.TransferOperation{}),
		c.RegisterType(&nftfx.Credential{}),
	)
	c.NextGroup()

	errs.Add(
		c.RegisterType(&propertyfx.MintOutput{}),
		c.RegisterType(&propertyfx.OwnedOutput{}),
		c.RegisterType(&propertyfx.MintOperation{}),
		c.RegisterType(&propertyfx.BurnOperation{}),
		c.RegisterType(&propertyfx.Credential{}),
	)

	return c
}
