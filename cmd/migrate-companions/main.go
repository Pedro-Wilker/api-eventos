// Backfill de companion_qr_codes + companion_checked_in para convidados
// existentes. Idempotente: pula convidados já preenchidos.
//
// Uso no servidor:
//   cd ~/apps/api-eventos
//   go run ./cmd/migrate-companions
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Pedro-Wilker/api-eventos/config"
	"github.com/Pedro-Wilker/api-eventos/models"
	"github.com/Pedro-Wilker/api-eventos/utils"
	"gorm.io/datatypes"
)

func toJSON(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}

func main() {
	config.ConnectDB()
	log.Println("== migrate-companions start ==")

	var guests []models.Guest
	if err := config.DB.Where("companion_qty > ?", 0).Find(&guests).Error; err != nil {
		log.Fatalf("query failed: %v", err)
	}
	log.Printf("convidados com acompanhantes: %d", len(guests))

	var updated, skipped, empty int
	for i := range guests {
		g := &guests[i]

		var existing []string
		if err := json.Unmarshal(g.CompanionQRCodes, &existing); err == nil && len(existing) > 0 {
			skipped++
			continue
		}

		var names []string
		if err := json.Unmarshal(g.CompanionNames, &names); err != nil || len(names) == 0 {
			empty++
			continue
		}

		codes := make([]string, len(names))
		checked := make([]bool, len(names))
		// Hash canonico sobre nome COMO SALVO em companion_names (whitespace
		// preservado). PDFs antigos foram gerados sem TrimSpace antes do FNV-1a
		// — TrimSpace aqui quebraria match com PDFs ja emitidos.
		for i, n := range names {
			if n == "" {
				continue
			}
			codes[i] = utils.GenerateCompanionQR(g.QRCode, n)
		}

		g.CompanionQRCodes = toJSON(codes)
		g.CompanionCheckedIn = toJSON(checked)

		if err := config.DB.Save(g).Error; err != nil {
			log.Printf("guest %d (%s) save failed: %v", g.ID, g.Name, err)
			continue
		}
		updated++
		log.Printf("[ok] guest %d (%s): %d acompanhante(s) -> %v", g.ID, g.Name, len(names), codes)
	}

	log.Printf("== migrate-companions done: atualizados=%d ignorados=%d sem_nomes=%d total=%d ==",
		updated, skipped, empty, len(guests))
	fmt.Printf("OK atualizados=%d ignorados=%d sem_nomes=%d total=%d\n", updated, skipped, empty, len(guests))
}
