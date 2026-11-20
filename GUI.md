# SpeedDisc GUI | Oficina portátil

A Windows x64 workbench for the existing SpeedDisc tool catalogue. It uses
Go and Windows controls, with no installer, web view, browser engine, internet
service or separately installed runtime. The legacy command-line edition
remains available.

## Usar / Use

- Escolha um pacote ou ferramenta, ou procure pelo nome. A aplicação inicia em
  **Só ver**, sem executar nada. A interface inclui PT/EN e tema claro/escuro.
- Leia a descrição e confirme qualquer alteração. Quando necessário, o Windows
  abre uma segunda instância como administrador, sem executar nada de imediato.
  Se indicar outra conta, essa instância usa o perfil dessa conta.
- No arranque, a análise é só de leitura. **Gerir entradas** permite rever uma
  entrada antes de a desactivar. O Winsock começa desligado; é uma escolha
  separada da limpeza da cache DNS.
- Após cada execução, incluindo uma execução com ressalvas, pode escolher
  **Exportar TXT** ou **Exportar PDF**. Pacotes têm um relatório com os resultados
  dos seus passos. Os ficheiros são locais e não substituem ficheiros existentes.
  Os relatórios podem conter informação do computador, como caminhos locais;
  reveja-os antes de partilhar. TXT mantém Unicode integral. O PDF usa o
  alfabeto coberto pela fonte e sinaliza glifos adicionais pelo código Unicode.
- `SpeedDisc-GUI.exe --demo` shows a fictional result without running maintenance.
  It allows report export, so the layout can be reviewed without using PC data.

## Build and verification

```sh
python scripts/build-gui.py --output dist/SpeedDisc-GUI.exe
python scripts/build-gui.py --verify dist/SpeedDisc-GUI.exe
go test ./...
```

The builder verifies the actual x64 Windows GUI subsystem, icon, asInvoker
manifest and Per-Monitor V2 DPI declaration. It generates its own resources
without changing the CLI. The Windows smoke workflow uses test fixtures, opens
and closes the app without maintenance, and tries to capture the regular and
fictional-report windows. Hosted Windows Server testing is not interactive
Windows 10/11 acceptance. Before promoting a public binary, verify a real
standard-user session (with and without UAC), window scale/keyboard navigation,
and one read-only estimate plus TXT/PDF save in an isolated environment. Test
changing operations only on disposable Windows fixtures.
