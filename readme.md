# Запуск проекта

## Локальный запуск

### Зависимости

- Go версии 1.27.3 или новее.
- Python 3 и пакет `matplotlib` — для построения изображения гистограммы.
- Исходное изображение в формате PNG.

Установите `matplotlib`, если он еще не установлен:

```sh
python -m pip install matplotlib
```

Запускайте команды из корня проекта. Для изображения по умолчанию (`imgs/gradient.png`) выполните:

```sh
go run ./src
python scripts/plot.py
```

Для другого изображения передайте путь к PNG аргументом:

```sh
go run ./src imgs/noise_texture.png
python scripts/plot.py
```

Go-программа создает зашифрованные и расшифрованные изображения в `imgs/`, CSV с метриками и данными гистограммы — в `results/`, а отчет — в `report.md`. Команда Python строит `results/histogram.png`, на который ссылается отчет.

Вместо аргумента можно задать исходное изображение через `SRC_IMG_PATH`. В PowerShell:

```powershell
$env:SRC_IMG_PATH = "imgs/noise_texture.png"
go run ./src
python scripts/plot.py
```

Для локального запуска переменные путей по умолчанию такие:

| Переменная окружения | Значение по умолчанию | Назначение |
| --- | --- | --- |
| `SRC_IMG_PATH` | `imgs/gradient.png` | Исходное PNG |
| `HISTOGRAM_CSV_PATH` | `results/histogram.csv` | CSV с данными гистограммы |
| `HISTOGRAM_IMG_PATH` | `results/histogram.png` | Изображение гистограммы |

`HISTOGRAM_CSV_PATH` и `HISTOGRAM_IMG_PATH` должны совпадать для Go-программы и `scripts/plot.py`, поскольку скрипт использует эти переменные для чтения CSV и сохранения изображения.

## Запуск в Docker

### Зависимости

- Docker и Docker Compose.
- Исходное изображение в формате PNG в каталоге `imgs/`.

### Команды

Запуск с изображением по умолчанию (`imgs/gradient.png`):

```sh
docker compose up --build
```

Запуск с другим изображением, например `imgs/photo.png`:

- **Linux/macOS:**

    ```sh
    SRC_IMG_PATH=/app/imgs/photo.png docker compose up --build
    ```

- **Windows (PowerShell):**

    ```PowerShell
    $env:SRC_IMG_PATH="/app/imgs/photo.png"; docker compose up --build
    ```

- **Универсально (файл .env рядом с docker-compose.yml)**

    ```env
    SRC_IMG_PATH=/app/imgs/photo.png
    ```

    И запуск стандартной командой
    ```sh
    docker compose up --build
    ```

### Параметры

| Переменная окружения | Значение по умолчанию | Назначение |
| --- | --- | --- |
| `SRC_IMG_PATH` | `/app/imgs/gradient.png` | Исходное PNG внутри контейнера |
| `HISTOGRAM_CSV_PATH` | `/app/results/histogram.csv` | Путь сохранения данных гистограммы |
| `HISTOGRAM_IMG_PATH` | `/app/results/histogram.png` | Путь сохранения изображения гистограммы |

Каталог `imgs/` и каталог `results/` доступны контейнеру через Docker Compose и сохраняют результаты на хосте.

### Пример шифрования и расшифрования

При запуске программа автоматически шифрует выбранное изображение, затем расшифровывает его и проверяет совпадение исходных и восстановленных данных. Для `photo.png` файлы появятся в `imgs/`:

- XOR: `photoXOR-enc.png` (шифрование), `photoXOR-dec.png` (расшифрование);
- AES-ECB: `photoECB-enc.png`, `photoECB-dec.png`;
- AES-CTR: `photoCTR-enc.png`, `photoCTR-dec.png`.

Отдельная команда для запуска только шифрования или только расшифрования не предусмотрена.

### Реализованные алгоритмы

- XOR с генератором Xorshift32;
- AES-128 в режимах ECB и CTR.
