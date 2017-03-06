package aci

import (
	"testing"
)

func TestVlanpool1(t *testing.T) {
	vlanspoolSplitTest(t, "", "", "")
	vlanspoolSplitTest(t, "vlanns-[a]-b", "a", "b")
	vlanspoolSplitTest(t, "vlanns-[]-", "", "")
	vlanspoolSplitTest(t, "vlanns-[a]-", "a", "")
	vlanspoolSplitTest(t, "vlanns-[]-b", "", "b")
	vlanspoolSplitTest(t, "vlanns-[a-b", "a", "b")
	vlanspoolSplitTest(t, "vlanns-a]-b", "a", "b")
	vlanspoolSplitTest(t, "vlanns-[-b", "", "b")
	vlanspoolSplitTest(t, "vlanns-]-b", "", "b")
	vlanspoolSplitTest(t, "vlanns-a-b", "a", "b")
	vlanspoolSplitTest(t, "vlanns--", "", "")
}

func vlanspoolSplitTest(t *testing.T, input, wantPool, wantMode string) {
	resultPool, resultMode := vlanpoolSplit(input)
	if resultPool != wantPool || resultMode != wantMode {
		t.Errorf("input=%s wantPool=%s gotPool=%s wantMode=%s gotMode=%s", input, wantPool, resultPool, wantMode, resultMode)
	}
}
