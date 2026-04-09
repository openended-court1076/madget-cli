# MadGet - GO Tabanlı Paket Yöneticisi

MadGet, Go diliyle geliştirilmiş, XML tabanlı yapılandırma dosyalarını kullanan örnek bir paket yönetim aracıdır. Bu proje, modern CLI araçlarının çalışma mantığını anlamak, Cobra kütüphanesini uygulamak ve yapılandırılmış verileri (XML) terminal ortamında işlemek amacıyla geliştirilmiştir.

## Özellikler
- Paket hakkında detaylı bilgi (adı, sürümü, açıklama)
- Kullanım kılavuzları ve örnek komutlar

## Kurulum
1. GitHub deposunu klonlayın: `git clone https://github.com/mehmetalidsy/madget-cli.git`
2. Proje dizinine geçin: `cd madget`
3. Gerekli bağımlılıkları indirin: `go mod tidy`
4. Çalıştırın: `go run cmd/madget.go` veya derleyin ve çalıştırın: `go build -o madget`

## Kullanım
- `madget info <paket>`: Paket hakkında bilgi gösterir.
- `madget install <paket>`: Paket yükler.
- Diğer komutlar için komut listedeki kullanım bilgisini inceleyin.

## Katkıda Bulundurma
Katkıda bulunmak isterseniz:
1. Depoyu çatallayın: `git fork`
2. Değişiklikler yapın ve commit edin: `git commit -am "Yaptığım değişiklikler"`
3. Çatalladığınız dallardan birinde pull request oluşturun.

## Lisans
Bu proje MIT Lisansı altında veröffentlichtir - https://github.com/mehmetalidsy/madget-cli/blob/main/LICENSE dosyasını görün.