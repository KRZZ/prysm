package blockchain

import (
	"testing"

	"github.com/prysmaticlabs/prysm/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/consensus-types/wrapper"
	enginev1 "github.com/prysmaticlabs/prysm/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/testing/require"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func Test_logStateTransitionData(t *testing.T) {
	payloadBlk := &ethpb.BeaconBlockBellatrix{
		Body: &ethpb.BeaconBlockBodyBellatrix{
			SyncAggregate: &ethpb.SyncAggregate{},
			ExecutionPayload: &enginev1.ExecutionPayload{
				BlockHash:    []byte{1, 2, 3},
				Transactions: [][]byte{{}, {}},
			},
		},
	}
	wrappedPayloadBlk, err := wrapper.WrappedBeaconBlock(payloadBlk)
	require.NoError(t, err)
	tests := []struct {
		name string
		b    func() interfaces.BeaconBlock
		want string
	}{
		{name: "empty block body",
			b: func() interfaces.BeaconBlock {
				wb, err := wrapper.WrappedBeaconBlock(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{}})
				require.NoError(t, err)
				return wb
			},
			want: "\"Finished applying state transition\" prefix=blockchain slot=0",
		},
		{name: "has attestation",
			b: func() interfaces.BeaconBlock {
				wb, err := wrapper.WrappedBeaconBlock(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{Attestations: []*ethpb.Attestation{{}}}})
				require.NoError(t, err)
				return wb
			},
			want: "\"Finished applying state transition\" attestations=1 prefix=blockchain slot=0",
		},
		{name: "has deposit",
			b: func() interfaces.BeaconBlock {
				wb, err := wrapper.WrappedBeaconBlock(
					&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{
						Attestations: []*ethpb.Attestation{{}},
						Deposits:     []*ethpb.Deposit{{}}}})
				require.NoError(t, err)
				return wb
			},
			want: "\"Finished applying state transition\" attestations=1 deposits=1 prefix=blockchain slot=0",
		},
		{name: "has attester slashing",
			b: func() interfaces.BeaconBlock {
				wb, err := wrapper.WrappedBeaconBlock(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{
					AttesterSlashings: []*ethpb.AttesterSlashing{{}}}})
				require.NoError(t, err)
				return wb
			},
			want: "\"Finished applying state transition\" attesterSlashings=1 prefix=blockchain slot=0",
		},
		{name: "has proposer slashing",
			b: func() interfaces.BeaconBlock {
				wb, err := wrapper.WrappedBeaconBlock(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{
					ProposerSlashings: []*ethpb.ProposerSlashing{{}}}})
				require.NoError(t, err)
				return wb
			},
			want: "\"Finished applying state transition\" prefix=blockchain proposerSlashings=1 slot=0",
		},
		{name: "has exit",
			b: func() interfaces.BeaconBlock {
				wb, err := wrapper.WrappedBeaconBlock(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{
					VoluntaryExits: []*ethpb.SignedVoluntaryExit{{}}}})
				require.NoError(t, err)
				return wb
			},
			want: "\"Finished applying state transition\" prefix=blockchain slot=0 voluntaryExits=1",
		},
		{name: "has everything",
			b: func() interfaces.BeaconBlock {
				wb, err := wrapper.WrappedBeaconBlock(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{
					Attestations:      []*ethpb.Attestation{{}},
					Deposits:          []*ethpb.Deposit{{}},
					AttesterSlashings: []*ethpb.AttesterSlashing{{}},
					ProposerSlashings: []*ethpb.ProposerSlashing{{}},
					VoluntaryExits:    []*ethpb.SignedVoluntaryExit{{}}}})
				require.NoError(t, err)
				return wb
			},
			want: "\"Finished applying state transition\" attestations=1 attesterSlashings=1 deposits=1 prefix=blockchain proposerSlashings=1 slot=0 voluntaryExits=1",
		},
		{name: "has payload",
			b:    func() interfaces.BeaconBlock { return wrappedPayloadBlk },
			want: "\"Finished applying state transition\" payloadHash=0x010203 prefix=blockchain slot=0 syncBitsCount=0 txCount=2",
		},
	}
	for _, tt := range tests {
		hook := logTest.NewGlobal()
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, logStateTransitionData(tt.b()))
			require.LogsContain(t, hook, tt.want)
		})
	}
}
