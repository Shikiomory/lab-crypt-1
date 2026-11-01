package main

import (
	"fmt"
	"os"
	"text/template"
)

type ReportData struct {
	Rows         []MetricRow
	NPCR         float64
	UACI         float64
	HistogramImg string
}

const reportTemplate = `#Отчет по лабораторной работе по криптографии

---

## 1. Постановка задачи
* **Реализованные методы шифрования:**
  * Потоковое шифрование XOR на базе ПСЧ Xorshift32;
  * Симметричное блочное шифрование AES-128 в режиме ECB;
  * Симметричное блочное шифрование AES-128 в режиме CTR.
* **Использованные изображения:** Синтетический растровый градиент, шахматная доска, шум.
* **Объекты сравнения:** Энтропия Шеннона, корреляция Пирсона (H/V), гистограммы, чувствительность к ключу (NPCR/UACI) и обратимость.

## 2. Реализация

* **XOR (Xorshift32):** Потоковое сложение байт изображения с псевдослучайной гаммой.
* **AES-ECB и AES-CTR:** Использован пакет crypto/aes со 128-битным ключом.
* **Вектор инициализации (IV):** Сгенерирован криптостойким генератором crypto/rand для режима CTR и сохранен в JSON-файл вместе с алгоритмом и исходным изображением.
* **Обратимость:** Подтверждена побайтовым совпадением до шифрования и после дешифрования.

## 3. Методика эксперимента
* **Исследуемое изображение:** Основное статистическое сравнение проводилось для изображения шахматной доски.
* **Рассчитываемые метрики:**
  * Информационная энтропия Шеннона ($H$) для оценки плотности распределения информации;
  * Коэффициент корреляции Пирсона между соседними отсчётами по горизонтали ($r_H$) и вертикали ($r_V$) для яркостной компоненты ($Y = 0.299R + 0.587G + 0.114B$);
  * 256-уровневые гистограммы распределения частот байт;
  * Показатели лавинного эффекта (NPCR и UACI);
* **Выборка пикселей для оценки корреляции:**
  * Анализировалась полная генеральная совокупность смежных пар пикселей изображения без прореживания;
* * **Эксперимент с изменением ключа:**
  * Исследовался режим AES-CTR на одном и том же открытом изображении;
  * Вектор инициализации (IV) зафиксирован неизменным для устранения сторонних факторов;
  * Сформированы два ключа $K_1$ и $K_2$, различающиеся ровно на 1 символ;
  * По полученным шифротекстам $C_1$ и $C_2$ рассчитаны метрики NPCR (процент изменившихся байтов) и UACI (средняя глубина различий).
## 4. Результаты

### Сводная таблица метрик
| Метод | Энтропия (бит/байт) | Корреляция H | Корреляция V |
| :--- | :---: | :---: | :---: |
{{- range .Rows }}
| **{{ .Method }}** | {{ printf "%.6f" .Entropy }} | {{ printf "%.6f" .CorrH }} | {{ printf "%.6f" .CorrV }} |
{{- end }}

### Чувствительность к изменению ключа

* **NPCR:** {{ printf "%.4f" .NPCR }}% 
* **UACI:** {{ printf "%.4f" .UACI }}%

### Частотное распределение байт (Гистограмма)

![Гистограммы распределения]({{.HistogramImg}})

## 5. Анализ
{{- $orig := index .Rows 0 -}}
{{- $xor  := index .Rows 1 -}}
{{- $ecb  := index .Rows 2 -}}
{{- $ctr  := index .Rows 3 }}
* **Энтропия и распределение байт:**
  * У исходного изображения энтропия минимальна ($H \approx {{ printf "%.2f" $orig.Entropy }}$), так как шахматная доска состоит всего из двух однородных цветов.
  * Для **XOR** ($H \approx {{ printf "%.2f" $xor.Entropy }}$) и **AES-CTR** ($H \approx {{ printf "%.2f" $ctr.Entropy }}$) энтропия вплотную приближается к максимуму ($8.0$), а частотное распределение байт на гистограмме равномерно.
  * Для **AES-ECB** энтропия падает до ${{ printf "%.2f" $ecb.Entropy }}$, а на графике видны пики: одинаковые 16-байтные блоки клеток шифруются в одинаковые блоки шифротекста, сохраняя избыточность.

* **Корреляция соседних пикселей:**
  * В оригинале зависимость между пикселями близка к единице ($r \approx {{ printf "%.2f" $orig.CorrH }}$).
  * **XOR** и **AES-CTR** полностью разрушают связи: коэффициенты падают практически до нуля ($r_H \approx {{ printf "%.4f" $ctr.CorrH }}$, $r_V \approx {{ printf "%.4f" $ctr.CorrV }}$).
  * **AES-ECB** демонстрирует аномалию: вертикальная корреляция ($r_V \approx {{ printf "%.2f" $ecb.CorrV }}$) даже превысила оригинал, а горизонтальная стала отрицательной ($r_H \approx {{ printf "%.2f" $ecb.CorrH }}$). Границы клеток и вертикальные полосы не зашифровались, а лишь изменили значения байт.

* **Причина отличия ECB от CTR:**
  * В **ECB** блоки открытого текста шифруются изолированно ($C_i = E_K(P_i)$), поэтому повторяющиеся клетки доски сохраняют визуальную форму.
  * В **CTR** шифруется монотонный счётчик с вектором инициализации ($C_i = P_i \oplus E_K(T_i)$), поэтому даже одинаковые блоки накладываются на уникальную псевдослучайную последовательность.

  * **Графики демонстрируют согласованность с численными показателями энтропии и корреляции.** 

## 6. Итоговый вывод
{{- $orig := index .Rows 0 -}}
{{- $xor  := index .Rows 1 -}}
{{- $ecb  := index .Rows 2 -}}
{{- $ctr  := index .Rows 3 }}
Эксперимент наглядно показал, что для защиты визуальных данных критически важен выбор режима работы шифра:

1. **AES-ECB статистически непригоден для изображений:** на структурах с повторяющимися элементами он сохраняет визуальные контуры, демонстрирует сильную остаточную корреляцию ($r_V \approx {{ printf "%.2f" $ecb.CorrV }}$) и резкую просадку энтропии ($H \approx {{ printf "%.2f" $ecb.Entropy }}$).
2. **XOR (Xorshift32)** формально обеспечивает равномерный шум ($H \approx {{ printf "%.2f" $xor.Entropy }}, r \approx 0$), однако не обладает криптостойкостью из-за линейности генератора и уязвимости к известному открытому тексту.
3. **AES-CTR является наиболее надёжным решением:** режим полностью ликвидирует пространственную корреляцию пикселей ($r \approx 0$), обеспечивает максимальную энтропию ($H \approx {{ printf "%.4f" $ctr.Entropy }}$), гарантирует строгую обратимость и показывает эталонный лавинный эффект ($\text{NPCR} \approx {{ printf "%.2f" .NPCR }}\%$, $\text{UACI} \approx {{ printf "%.2f" .UACI }}\%$).
`

func ClearFile(filePath string) error {
	return os.Truncate(filePath, 0)
}

func CreateTable(filePath string, rows []MetricRow) error {
	outFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer outFile.Close()

	header := "| Method | Entropy | CorrH | CorrV |\n"
	separator := "|--------|---------|-------|-------|\n"
	_, err = outFile.WriteString(header)
	if err != nil {
		return fmt.Errorf("не удалось записать заголовок: %w", err)
	}
	_, err = outFile.WriteString(separator)
	if err != nil {
		return fmt.Errorf("не удалось записать разделитель: %w", err)
	}

	for _, row := range rows {
		line := fmt.Sprintf("| %s | %f | %f | %f |\n", row.Method, row.Entropy, row.CorrH, row.CorrV)
		_, err = outFile.WriteString(line)
		if err != nil {
			return fmt.Errorf("Не удалось записать строку: %w", err)
		}
	}

	_, err = outFile.WriteString("\n")
	if err != nil {
		return fmt.Errorf("Не удалось записать строку: %w", err)
	}
	return nil
}

func WriteKeyComparison(filePath string, npcr, uaci float64) error {
	outFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Не удалось открыть файл: %w", err)
	}

	defer outFile.Close()

	_, err = outFile.WriteString("Чувствительность к изменению ключа:\n")
	if err != nil {
		return fmt.Errorf("Не удалось записать заголовок: %w", err)
	}
	_, err = outFile.WriteString(fmt.Sprintf("NPCR: %f\n", npcr))
	if err != nil {
		return fmt.Errorf("Не удалось записать NPCR: %w", err)
	}

	_, err = outFile.WriteString(fmt.Sprintf("UACI: %f\n", uaci))
	if err != nil {
		return fmt.Errorf("Не удалось записать UACI: %w", err)
	}
	return nil
}

func WriteHistogram(filePath string, imgFilePath string) error {
	outFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Не удалось открыть файл: %w", err)
	}
	defer outFile.Close()
	_, err = outFile.WriteString(fmt.Sprintf("Гистограмма:\n![гистограмма](%s)\n", imgFilePath))
	if err != nil {
		return fmt.Errorf("Не удалось записать гистограмму: %w", err)
	}
	return nil
}

func WriteReport(filePath string, rows []MetricRow, npcr, uaci float64, imgFilePath string) error {
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return fmt.Errorf("Не удалось создать отчет: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Не удалось создать файл отчета: %w", err)
	}
	defer file.Close()

	data := ReportData{
		Rows:         rows,
		NPCR:         npcr * 100,
		UACI:         uaci * 100,
		HistogramImg: imgFilePath,
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("Не удалось записать отчет: %w", err)
	}

	return nil
}
