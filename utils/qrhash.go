package utils

import (
	"fmt"
	"hash/fnv"
	"strings"
)

// FNV1a32Hex calcula hash FNV-1a 32 bits e devolve 8 chars hex lowercase.
// Implementação DEVE ser byte-idêntica a src/lib/qr.ts:fnv1a32 (offset 0x811c9dc5,
// prime 0x01000193, unsigned 32-bit).
func FNV1a32Hex(s string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return fmt.Sprintf("%08x", h.Sum32())
}

// GenerateCompanionQR devolve qr_code único e determinístico para um
// acompanhante: substitui o 1º segmento (8 chars) do qr do titular pelo
// hash FNV-1a de (titularQR + "|" + nomeAcomp). Idêntico a
// src/lib/qr.ts:gerarQrUnicoAcompanhante — IMPORTANTÍSSIMO manter sync.
func GenerateCompanionQR(titularQR, nomeAcomp string) string {
	partes := strings.Split(titularQR, "-")
	if len(partes) != 5 {
		return titularQR // fallback seguro
	}
	novoId := FNV1a32Hex(titularQR + "|" + nomeAcomp)
	return strings.Join([]string{novoId, partes[1], partes[2], partes[3], partes[4]}, "-")
}
