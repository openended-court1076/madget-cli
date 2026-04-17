# Örnek paket

- `MadGet.xml` — publish için manifest (`package_name`: **example-pkg**, `version`: **1.0.0**).
- `payload/` — tarball içine girecek dosyalar (kökte `README.txt`).

Tarball üretmek (repo kökünden):

```bash
make example-tgz
```

Windows PowerShell:

```powershell
.\scripts\build-example-package.ps1
```

Çıktı: `example/package.tgz` (`.gitignore` ile repoya eklenmez; yerelde üretilir.)

Ardından registry açıkken:

```bash
go run . publish ./example/MadGet.xml ./example/package.tgz
go run . install example-pkg@^1.0.0
```
