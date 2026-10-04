package main

import "fmt"

// text[key][0] is European Portuguese, text[key][1] is English.
var text = map[string][2]string{
	"subtitle": {
		"Manutenção portátil para técnicos de informática",
		"Portable maintenance for IT technicians",
	},
	"lang_now": {
		"Idioma: português",
		"Language: English",
	},
	"section_packages": {"Pacotes", "Packages"},
	"section_actions":  {"Ações", "Actions"},
	"opt_quick":        {"Pacote Rápido", "Quick package"},
	"desc_quick": {
		"Temporários, cache de transferências do Windows Update, miniaturas e plano Alto desempenho.",
		"Temp files, Windows Update download cache, thumbnail cache, and the High performance power plan.",
	},
	"opt_deep": {"Pacote Profundo", "Deep package"},
	"desc_deep": {
		"O conjunto rápido, mais cache do Delivery Optimization e relatórios de erros, SFC e DISM. Cria um ponto de restauro se o Windows deixar. Não desativa o arranque e não esvazia a Reciclagem.",
		"The quick set, plus Delivery Optimization cache and Windows error reports, SFC and DISM. Creates a restore point when Windows allows it. Does not disable startup items and does not empty the Recycle Bin.",
	},
	"opt_full": {"Pacote completo", "Full package"},
	"desc_full": {
		"Temporários, cache de transferências do Windows Update, miniaturas, cache do Delivery Optimization, relatórios de erros, plano Alto desempenho, SFC e DISM. Cria primeiro um ponto de restauro, se o Windows deixar. Não desativa o arranque e não esvazia a Reciclagem.",
		"Temp files, the Windows Update download cache, thumbnails, Delivery Optimization cache, error reports, the High performance power plan, SFC and DISM. Creates a restore point first when Windows allows it. Does not disable startup items and does not empty the Recycle Bin.",
	},
	"opt_look": {"Pacote Só ver", "Look only"},
	"desc_look": {
		"Mostra sistema, processador, memória, discos, arranque e uma estimativa dos temporários. Não altera nada.",
		"Shows OS, processor, memory, disks, startup count, and a temp-size estimate. Changes nothing.",
	},
	"opt_temp": {"Ficheiros temporários", "Temporary files"},
	"desc_temp": {
		"Apaga ficheiros que não estejam em uso nas pastas Temp do utilizador e do Windows.",
		"Deletes files that are not in use in the user Temp folder and in Windows Temp.",
	},
	"opt_wu": {"Cache do Windows Update", "Windows Update download cache"},
	"desc_wu": {
		"Apaga só o conteúdo de SoftwareDistribution\\Download. O resto dessa pasta fica intacto.",
		"Deletes only the contents of SoftwareDistribution\\Download. The rest of that tree is left alone.",
	},
	"opt_thumb": {"Cache de miniaturas", "Thumbnail cache"},
	"desc_thumb": {
		"Apaga thumbcache e iconcache se não estiverem em uso. Não fecha o Explorador.",
		"Deletes thumbcache and iconcache files that are not in use. Does not close Explorer.",
	},
	"opt_do": {"Cache do Delivery Optimization", "Delivery Optimization cache"},
	"desc_do": {
		"Apaga a cache do Delivery Optimization, se a pasta existir.",
		"Deletes the Delivery Optimization cache when that folder exists.",
	},
	"opt_power": {"Plano Alto desempenho", "High performance power plan"},
	"desc_power": {
		"Ativa o plano Alto desempenho com powercfg. Não altera serviços.",
		"Activates the High performance plan with powercfg. Does not change services.",
	},
	"opt_wer": {"Relatórios de erros", "Windows error reports"},
	"desc_wer": {
		"Apaga relatórios de erros já guardados pelo Windows (WER). Não mexe no registo.",
		"Deletes stored Windows Error Reporting files. Does not edit the registry.",
	},
	"opt_sfc": {"SFC /scannow", "SFC /scannow"},
	"desc_sfc": {
		"Corre sfc /scannow. Pode demorar muito. Não indica uma percentagem de velocidade.",
		"Runs sfc /scannow. It can take a long time. It does not report a speed percentage.",
	},
	"opt_dism": {"DISM RestoreHealth", "DISM RestoreHealth"},
	"desc_dism": {
		"Corre DISM /Online /Cleanup-Image /RestoreHealth. Pode demorar muito.",
		"Runs DISM /Online /Cleanup-Image /RestoreHealth. It can take a long time.",
	},
	"opt_restore": {"Ponto de restauro", "System restore point"},
	"desc_restore": {
		"Cria um ponto de restauro. Se o Windows recusar, fica no relatório e o resto segue.",
		"Creates a restore point. If Windows refuses, that is recorded and the rest continues.",
	},
	"opt_startup": {"Arranque (escolher)", "Startup (pick)"},
	"desc_startup": {
		"Lista chaves Run (HKLM e HKCU) e pastas de arranque. Só desativa o que escolher.",
		"Lists Run keys (HKLM and HKCU) and Startup folders. Disables only what you pick.",
	},
	"opt_info": {"Informação do sistema", "System information"},
	"desc_info": {
		"Igual ao pacote Só ver: lê informação e não altera nada.",
		"Same as Look only: reads information and changes nothing.",
	},
	"opt_lang": {"English", "Português"},
	"opt_quit": {"Sair", "Quit"},
	"menu_hint": {
		"Várias de uma vez, por exemplo 5, 9, 14",
		"Several at once, for example 5, 9, 14",
	},
	"bad_selection": {
		"Escolha inválida. Use números do menu, separados por vírgulas. Exemplo: 1 ou 5, 9, 14.",
		"Invalid choice. Use menu numbers, separated by commas. Example: 1 or 5, 9, 14.",
	},
	"need_admin": {
		"O SpeedDisc tem de correr como administrador. O Windows deve pedir elevação antes deste programa abrir. Se isso não aconteceu, feche e use Executar como administrador.",
		"SpeedDisc must run as administrator. Windows should ask for elevation before this program opens. If it did not, close it and use Run as administrator.",
	},
	"confirm_intro": {
		"Isto vai alterar o computador.",
		"This will change the computer.",
	},
	"confirm_long": {
		"O SFC e o DISM podem demorar muito tempo.",
		"SFC and DISM can take a long time.",
	},
	"confirm_restore": {
		"Será criado um ponto de restauro se o Windows o permitir. Se falhar, o resto continua e fica escrito no relatório.",
		"A restore point is created when Windows allows it. If that fails, the rest continues and the failure is written in the report.",
	},
	"confirm_startup": {
		"Nenhuma entrada de arranque é desativada sem uma segunda confirmação.",
		"No startup entry is disabled without a second confirmation.",
	},
	"confirm_safe": {
		"A Reciclagem não é esvaziada. Nenhum serviço é desativado. Os ficheiros em uso são ignorados.",
		"The Recycle Bin is not emptied. No service is disabled. Files that are in use are skipped.",
	},
	"confirm_ask": {
		"Continuar? (s/n) ",
		"Continue? (y/n) ",
	},
	"cancelled": {
		"Cancelado. Nada foi alterado nesta passagem.",
		"Cancelled. Nothing was changed this time.",
	},
	"press_enter": {
		"Enter para voltar ao menu...",
		"Enter to return to the menu...",
	},
	"press_enter_exit": {
		"Enter para fechar...",
		"Enter to close...",
	},
	"bye":     {"Até já.", "Done."},
	"running": {"A executar: %s", "Running: %s"},
	"long_running": {
		"Isto pode demorar.",
		"This can take a long time.",
	},
	"done_sfc_line":  {"SFC concluído", "SFC finished"},
	"done_dism_line": {"DISM concluído", "DISM finished"},
	"report_at":      {"Relatório: %s", "Report: %s"},
	"report_fail":    {"Não foi possível escrever o relatório: %s", "Could not write the report: %s"},
	"startup_empty": {
		"Não foram encontradas entradas de arranque.",
		"No startup entries were found.",
	},
	"startup_pick": {
		"Números a desativar, separados por vírgulas. Enter para não desativar nada.",
		"Numbers to disable, separated by commas. Enter to disable nothing.",
	},
	"startup_bad": {
		"Números inválidos. Use apenas números da lista.",
		"Invalid numbers. Use only numbers from the list.",
	},
	"startup_confirm": {
		"Desativar as entradas escolhidas? (s/n) ",
		"Disable the chosen entries? (y/n) ",
	},
	"rep_title":            {"SpeedDisc — relatório", "SpeedDisc — report"},
	"rep_version":          {"Versão: %s", "Version: %s"},
	"rep_date":             {"Data (hora local do PC): %s", "Date (PC local time): %s"},
	"rep_lang":             {"Idioma: %s", "Language: %s"},
	"rep_computer":         {"Computador: %s", "Computer: %s"},
	"rep_found":            {"Encontrado", "Found"},
	"rep_done":             {"Feito", "Done"},
	"rep_failed":           {"Falhou", "Failed"},
	"rep_notes":            {"Notas", "Notes"},
	"rep_none":             {"Nada a registar.", "Nothing to record."},
	"rep_lang_name":        {"português", "English"},
	"found_os":             {"Sistema: %s", "System: %s"},
	"found_cpu":            {"Processador: %s", "Processor: %s"},
	"found_ram":            {"Memória: %s livres de %s", "Memory: %s free of %s"},
	"found_disk":           {"Disco %s: %s livres de %s", "Disk %s: %s free of %s"},
	"found_startup_count":  {"Entradas de arranque: %d", "Startup entries: %d"},
	"found_startup_sample": {"Amostra de arranque: %s", "Startup sample: %s"},
	"found_temp": {
		"Temporários (estimativa, pode falhar o que não estiver acessível): %s",
		"Temporary files (estimate; inaccessible items may be missing): %s",
	},
	"done_clean": {
		"%s: %d ficheiros removidos (%s); %d ignorados (em uso ou sem acesso).",
		"%s: removed %d files (%s); skipped %d (in use or not accessible).",
	},
	"done_missing": {"%s: a pasta não existe.", "%s: the folder does not exist."},
	"done_power":   {"Plano de energia Alto desempenho ativado (%s).", "High performance power plan is active (%s)."},
	"done_restore": {"Ponto de restauro criado.", "Restore point created."},
	"done_cmd":     {"%s terminou com o código %d.", "%s finished with code %d."},
	"done_cmd_reboot": {
		"%s terminou com o código %d. O Windows pode pedir para reiniciar.",
		"%s finished with code %d. Windows may ask for a restart.",
	},
	"done_startup_reg": {
		"Arranque desativado: %s (%s). Valor anterior: %s",
		"Startup disabled: %s (%s). Previous value: %s",
	},
	"done_startup_file": {
		"Arranque desativado: %s. Movido para: %s",
		"Startup disabled: %s. Moved to: %s",
	},
	"done_startup_none": {
		"Arranque: nenhuma entrada foi desativada.",
		"Startup: no entry was disabled.",
	},
	"fail_clean": {"%s: %s", "%s: %s"},
	"fail_power": {
		"Não foi possível ativar o plano Alto desempenho. %s",
		"Could not activate the High performance plan. %s",
	},
	"fail_restore": {
		"Ponto de restauro não criado (a continuar). %s",
		"Restore point was not created (continuing). %s",
	},
	"fail_cmd":     {"%s falhou com o código %d. %s", "%s failed with code %d. %s"},
	"fail_startup": {"Não foi possível desativar %s. %s", "Could not disable %s. %s"},
	"fail_startup_type": {
		"Não foi possível desativar %s: o valor não é texto e não foi alterado.",
		"Could not disable %s: the value is not text and was left unchanged.",
	},
	"fail_startup_changed": {
		"Não foi possível desativar %s: o valor mudou desde a listagem e não foi alterado.",
		"Could not disable %s: the value changed after it was listed and was left unchanged.",
	},
	"fail_info":      {"%s: %s", "%s: %s"},
	"note_recycle":   {"A Reciclagem não foi esvaziada.", "The Recycle Bin was not emptied."},
	"note_services":  {"Nenhum serviço do Windows foi desativado.", "No Windows service was disabled."},
	"note_prefetch":  {"A pasta Prefetch não foi apagada.", "The Prefetch folder was not deleted."},
	"note_sysmain":   {"O SysMain não foi alterado.", "SysMain was not changed."},
	"name_user_temp": {"Temporários do utilizador", "User temp files"},
	"name_win_temp":  {"Temporários do Windows", "Windows temp files"},
	"name_wu":        {"Cache de transferências do Windows Update", "Windows Update download cache"},
	"name_thumb":     {"Cache de miniaturas", "Thumbnail cache"},
	"name_do":        {"Cache do Delivery Optimization", "Delivery Optimization cache"},
	"name_wer":       {"Relatórios de erros do Windows", "Windows error reports"},
	"name_power":     {"Plano de energia", "Power plan"},
	"name_sfc":       {"SFC /scannow", "SFC /scannow"},
	"name_dism":      {"DISM RestoreHealth", "DISM RestoreHealth"},
	"name_restore":   {"Ponto de restauro", "System restore point"},
	"name_startup":   {"Arranque", "Startup"},
	"name_os":        {"Sistema", "System"},
	"name_cpu":       {"Processador", "Processor"},
	"name_ram":       {"Memória", "Memory"},
	"name_temp":      {"Temporários", "Temporary files"},
	"name_disks":     {"Discos", "Disks"},
	"not_text_value": {"(valor não texto)", "(not a text value)"},
	"err_shallow": {
		"caminho demasiado curto; nada foi apagado",
		"path is too shallow; nothing was deleted",
	},
	"err_danger": {
		"caminho protegido do Windows; nada foi apagado",
		"protected Windows path; nothing was deleted",
	},
	"err_not_allowed": {
		"caminho fora da lista permitida; nada foi apagado",
		"path is not on the allowlist; nothing was deleted",
	},
	"err_symlink": {
		"atalho simbólico; não foi seguido",
		"symlink; it was not followed",
	},
	"err_changed": {
		"o valor mudou desde a listagem",
		"the value changed after it was listed",
	},
	"err_not_text": {
		"o valor não é texto",
		"the value is not text",
	},
	"err_powercfg": {
		"o powercfg não conseguiu ativar o plano",
		"powercfg could not activate the plan",
	},
	"err_scheme": {
		"o plano não ficou ativo",
		"the plan did not become active",
	},
	"err_no_disk": {
		"nenhum disco fixo",
		"no fixed disk",
	},
	"err_memory": {
		"não foi possível ler a memória",
		"could not read memory",
	},
	"err_exit": {
		"sem mais detalhe",
		"no further detail",
	},
}

type Lang int

const (
	LangPT Lang = 0
	LangEN Lang = 1
)

func otherLang(l Lang) Lang {
	if l == LangPT {
		return LangEN
	}
	return LangPT
}

func T(lang Lang, key string, args ...any) string {
	pair, ok := text[key]
	s := key
	if ok {
		s = pair[lang]
	}
	if len(args) == 0 {
		return s
	}
	return fmt.Sprintf(s, args...)
}
