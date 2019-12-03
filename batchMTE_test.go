// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ach

import (
	"log"
	"testing"
	"time"

	"github.com/ourly/base"
)

// mockBatchMTEHeader creates a MTE batch header
func mockBatchMTEHeader() *BatchHeader {
	bh := NewBatchHeader()
	bh.ServiceClassCode = DebitsOnly
	bh.CompanyName = "Merchant with ATM"
	bh.CompanyIdentification = "231380104"
	bh.StandardEntryClassCode = MTE
	bh.CompanyEntryDescription = "CASH WITHDRAW"
	bh.EffectiveEntryDate = time.Now().AddDate(0, 0, 1).Format("060102") // YYMMDD
	bh.ODFIIdentification = "23138010"
	return bh
}

// mockMTEEntryDetail creates a MTE entry detail
func mockMTEEntryDetail() *EntryDetail {
	entry := NewEntryDetail()
	entry.TransactionCode = CheckingDebit
	entry.SetRDFI("031300012")
	entry.DFIAccountNumber = "744-5678-99"
	entry.Amount = 10000
	entry.SetOriginalTraceNumber("031300010000001")
	entry.SetReceivingCompany("JANE DOE")
	entry.SetTraceNumber("23138010", 1)
	entry.AddendaRecordIndicator = 1

	addenda02 := NewAddenda02()

	// NACHA rules example: 200509*321 East Market Street*Anytown*VA\
	addenda02.TerminalIdentificationCode = "200509"
	addenda02.TerminalLocation = "321 East Market Street"
	addenda02.TerminalCity = "ANYTOWN"
	addenda02.TerminalState = "VA"

	addenda02.TransactionSerialNumber = "123456" // Generated by Terminal, used for audits
	addenda02.TransactionDate = "1224"
	addenda02.TraceNumber = entry.TraceNumber
	entry.Addenda02 = addenda02

	return entry
}

// mockBatchMTE creates a MTE batch
func mockBatchMTE() *BatchMTE {
	batch := NewBatchMTE(mockBatchMTEHeader())
	batch.AddEntry(mockMTEEntryDetail())
	if err := batch.Create(); err != nil {
		log.Fatalf("Unexpected error building batch: %s\n", err)
	}
	return batch
}

// testBatchMTEHeader creates a MTE batch header
func testBatchMTEHeader(t testing.TB) {
	batch, _ := NewBatch(mockBatchMTEHeader())
	_, ok := batch.(*BatchMTE)
	if !ok {
		t.Error("Expecting BatchMTE")
	}
}

// TestBatchMTEHeader tests creating a MTE batch header
func TestBatchMTEHeader(t *testing.T) {
	testBatchMTEHeader(t)
}

// BenchmarkBatchMTEHeader benchmark creating a MTE batch header
func BenchmarkBatchMTEHeader(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEHeader(b)
	}
}

// TestBatchMTEAddendum02 validates Addenda02 returns an error
func TestBatchMTEAddendum02(t *testing.T) {
	mockBatch := NewBatchMTE(mockBatchMTEHeader())
	mockBatch.AddEntry(mockMTEEntryDetail())
	err := mockBatch.Create()
	// TODO: are we expecting there to be an error here?
	if !base.Match(err, nil) {
		t.Errorf("%T: %s", err, err)
	}
}

// testBatchMTEReceivingCompanyName validates Receiving company / Individual name is a mandatory field
func testBatchMTEReceivingCompanyName(t testing.TB) {
	mockBatch := mockBatchMTE()
	// modify the Individual name / receiving company to nothing
	mockBatch.GetEntries()[0].SetReceivingCompany("")
	err := mockBatch.Validate()
	if !base.Match(err, ErrConstructor) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchMTEReceivingCompanyName tests validating receiving company / Individual name is a mandatory field
func TestBatchMTEReceivingCompanyName(t *testing.T) {
	testBatchMTEReceivingCompanyName(t)
}

// BenchmarkBatchMTEReceivingCompanyName benchmarks validating receiving company / Individual name is a mandatory field
func BenchmarkBatchMTEReceivingCompanyName(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEReceivingCompanyName(b)
	}
}

// testBatchMTEAddendaTypeCode validates addenda type code is 05
func testBatchMTEAddendaTypeCode(t testing.TB) {
	mockBatch := mockBatchMTE()
	mockBatch.GetEntries()[0].Addenda02.TypeCode = "05"
	err := mockBatch.Validate()
	if !base.Match(err, ErrAddendaTypeCode) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchMTEAddendaTypeCode tests validating addenda type code is 05
func TestBatchMTEAddendaTypeCode(t *testing.T) {
	testBatchMTEAddendaTypeCode(t)
}

// BenchmarkBatchMTEAddendaTypeCod benchmarks validating addenda type code is 05
func BenchmarkBatchMTEAddendaTypeCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEAddendaTypeCode(b)
	}
}

// testBatchMTESEC validates that the standard entry class code is MTE for batchMTE
func testBatchMTESEC(t testing.TB) {
	mockBatch := mockBatchMTE()
	mockBatch.Header.StandardEntryClassCode = ACK
	err := mockBatch.Validate()
	if !base.Match(err, ErrBatchSECType) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchMTESEC tests validating that the standard entry class code is MTE for batchMTE
func TestBatchMTESEC(t *testing.T) {
	testBatchMTESEC(t)
}

// BenchmarkBatchMTESEC benchmarks validating that the standard entry class code is MTE for batch MTE
func BenchmarkBatchMTESEC(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTESEC(b)
	}
}

// testBatchMTEServiceClassCode validates ServiceClassCode
func testBatchMTEServiceClassCode(t testing.TB) {
	mockBatch := mockBatchMTE()
	// Batch Header information is required to Create a batch.
	mockBatch.GetHeader().ServiceClassCode = 0
	err := mockBatch.Create()
	if !base.Match(err, ErrConstructor) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchMTEServiceClassCode tests validating ServiceClassCode
func TestBatchMTEServiceClassCode(t *testing.T) {
	testBatchMTEServiceClassCode(t)
}

// BenchmarkBatchMTEServiceClassCode benchmarks validating ServiceClassCode
func BenchmarkBatchMTEServiceClassCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEServiceClassCode(b)
	}
}

// TestBatchMTEAmount validates Amount
func TestBatchMTEAmount(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.GetEntries()[0].Amount = 0
	err := mockBatch.Create()
	if !base.Match(err, ErrBatchAmountZero) {
		t.Errorf("%T: %s", err, err)
	}
}

func TestBatchMTETerminalState(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.GetEntries()[0].Addenda02.TerminalState = "XX"
	err := mockBatch.Create()
	if !base.Match(err, ErrValidState) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchMTEIndividualName validates IndividualName
func TestBatchMTEIndividualName(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.GetEntries()[0].IndividualName = ""
	err := mockBatch.Create()
	if !base.Match(err, ErrConstructor) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchMTEIdentificationNumber validates IdentificationNumber
func TestBatchMTEIdentificationNumber(t *testing.T) {
	mockBatch := mockBatchMTE()

	// NACHA rules state MTE records can't be all spaces or all zeros
	mockBatch.GetEntries()[0].IdentificationNumber = "   "
	err := mockBatch.Validate()
	if !base.Match(err, ErrIdentificationNumber) {
		t.Errorf("%T: %s", err, err)
	}

	mockBatch.GetEntries()[0].IdentificationNumber = "000000"
	err = mockBatch.Validate()
	if !base.Match(err, ErrIdentificationNumber) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchMTEValidTranCodeForServiceClassCode validates a transactionCode based on ServiceClassCode
func TestBatchMTEValidTranCodeForServiceClassCode(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.GetHeader().ServiceClassCode = CreditsOnly
	err := mockBatch.Create()
	if !base.Match(err, NewErrBatchServiceClassTranCode(CreditsOnly, 27)) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchMTEAddenda05 validates BatchMTE cannot have Addenda05
func TestBatchMTEAddenda05(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.Entries[0].AddendaRecordIndicator = 1
	mockBatch.GetEntries()[0].AddAddenda05(mockAddenda05())
	err := mockBatch.Create()
	if !base.Match(err, ErrBatchAddendaCategory) {
		t.Errorf("%T: %s", err, err)
	}
}
