import csv
import os
import sys
import matplotlib.pyplot as plt

def main():
    csv_path = os.getenv('HISTOGRAM_CSV_PATH', 'results/histogram.csv')
    out_path = os.getenv('HISTOGRAM_IMG_PATH', 'results/histogram.png')

    if not os.path.exists(csv_path):
        print(f'Файл {csv_path} не найден. Пожалуйста, убедитесь, что он существует.')
        sys.exit(1)

    bytes_range = []
    orig = []
    xor = []
    ecb = []
    ctr = []

    with open(csv_path, 'r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            bytes_range.append(int(row['byte']))
            orig.append(int(row['original']))
            xor.append(int(row['xor']))
            ecb.append(int(row['ecb']))
            ctr.append(int(row['ctr']))

    fig, axes = plt.subplots(2, 2,figsize=(12, 6))
    configs = [
        (axes[0,0], orig, 'Оригинал', 'blue'),
        (axes[0,1], xor, 'XOR', 'orange'),
        (axes[1,0], ecb, 'ECB', 'green'),
        (axes[1,1], ctr, 'CTR', 'red')
    ]

    for ax, data, title, color in configs:
        ax.bar(bytes_range, data, width=1.0, color=color, edgecolor='none')
        ax.set_title(title)
        ax.set_xlabel('Байт')
        ax.set_ylabel('Количество')
        ax.set_xlim(0, 255)
        ax.grid(True, linestyle='--', alpha=0.5)

    plt.tight_layout()
    plt.savefig(out_path, dpi=300)
    print(f'\nГистограмма сохранена в {out_path}')

if __name__ == '__main__':
    main()