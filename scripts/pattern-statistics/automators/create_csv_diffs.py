import os
import seaborn as sns
import matplotlib.pyplot as plt
import pandas as pd
from matplotlib.colors import LinearSegmentedColormap


def generate_heatmaps(root_dir):
    # Base directory path
    base_dir = os.path.join(root_dir, './csv/balanced')

    pattern_files = []
    random_files = []

    if not os.path.exists(os.path.join(root_dir, './csv/diff')):
        os.mkdir(os.path.join(root_dir, './csv/diff'))

    for root, _, files in os.walk(base_dir):
        for file in files:
            if file[-len("random.csv"):] == "random.csv":
                random_files.append(file)
            else:
                pattern_files.append(file)

    for pattern_file in pattern_files:
        params = pattern_file.split("_")
        random_file = ""
        for rfile in random_files:
            if rfile[:len("_".join(params[:2]))] == "_".join(params[:2]):
                random_file = rfile
                break

        df_random = pd.read_csv(os.path.join(root_dir, './csv/balanced', random_file), delimiter=',', index_col=0)
        df_pattern = pd.read_csv(os.path.join(root_dir, './csv/balanced', pattern_file), delimiter=',', index_col=0)

        df_diff = df_pattern - df_random

        print(os.path.join(root_dir, './csv/diff', pattern_file))
        df_diff.to_csv(os.path.join(root_dir, './csv/diff', pattern_file))


if __name__ == '__main__':
    generate_heatmaps("/Users/ruae/Projects/strategy-evaluator/output_7")
    generate_heatmaps("/Users/ruae/Projects/strategy-evaluator/output_14")
    generate_heatmaps("/Users/ruae/Projects/strategy-evaluator/output_21")
