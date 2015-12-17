// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"strings"
	"testing"

	"github.com/keybase/client/go/libkb"
)

func TestSaltPackEncrypt(t *testing.T) {
	tc := SetupEngineTest(t, "SaltPackEncrypt")
	defer tc.Cleanup()

	u1 := CreateAndSignupFakeUser(tc, "nalcp")
	u2 := CreateAndSignupFakeUser(tc, "nalcp")
	u3 := CreateAndSignupFakeUser(tc, "nalcp")

	trackUI := &FakeIdentifyUI{
		Proofs: make(map[string]string),
	}
	ctx := &Context{IdentifyUI: trackUI, SecretUI: u3.NewSecretUI()}

	run := func(Recips []string) {
		sink := libkb.NewBufferCloser()
		arg := &SaltPackEncryptArg{
			Recips: Recips,
			Source: strings.NewReader("id2 and encrypt, id2 and encrypt"),
			Sink:   sink,
		}

		eng := NewSaltPackEncrypt(arg, tc.G)
		if err := RunEngine(eng, ctx); err != nil {
			t.Fatal(err)
		}

		out := sink.Bytes()
		if len(out) == 0 {
			t.Fatal("no output")
		}
	}
	run([]string{u1.Username, u2.Username})

	// If we add ourselves, we should be smart and not error out
	// (We are u3 in this case)
	run([]string{u1.Username, u2.Username, u3.Username})
}

func TestSaltPackEncryptSelfNoKey(t *testing.T) {
	tc := SetupEngineTest(t, "SaltPackEncrypt")
	defer tc.Cleanup()

	_, passphrase := createFakeUserWithNoKeys(tc)
	trackUI := &FakeIdentifyUI{
		Proofs: make(map[string]string),
	}
	ctx := &Context{IdentifyUI: trackUI, SecretUI: &libkb.TestSecretUI{Passphrase: passphrase}}

	sink := libkb.NewBufferCloser()
	arg := &SaltPackEncryptArg{
		Recips: []string{"t_alice", "t_bob+kbtester1@twitter", "t_charlie+tacovontaco@twitter"},
		Source: strings.NewReader("track and encrypt, track and encrypt"),
		Sink:   sink,
	}

	eng := NewSaltPackEncrypt(arg, tc.G)
	err := RunEngine(eng, ctx)
	if _, ok := err.(libkb.NoKeyError); !ok {
		t.Fatalf("expected error type libkb.NoKeyError, got %T (%s)", err, err)
	}
}
