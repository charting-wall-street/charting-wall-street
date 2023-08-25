import os
import seaborn as sns
import matplotlib.pyplot as plt
import pandas as pd
from matplotlib.colors import LinearSegmentedColormap

# Define configurations based on directory type
configurations = {
    'balanced': {
        'cmap_center': 50,
        'vmin': 25,
        'vmax': 75
    },
    'worst': {
        'cmap_center': 50,
        'vmin': 25,
        'vmax': 75
    },
    'diff': {
        'cmap_center': 0,
        'vmin': -10,
        'vmax': 10
    },
    'size': {
        'cmap_center': 5000,
        'vmin': 0,
        'vmax': 50000
    }
}


def generate_heatmaps(root_dir):
    # Base directory path
    base_dir = os.path.join(root_dir, './csv/')

    # Create SVG directory if it doesn't exist
    if not os.path.exists(os.path.join(root_dir, './svg')):
        os.mkdir(os.path.join(root_dir, './svg'))

    # Loop through directories and files
    for root, _, files in os.walk(base_dir):
        for file in files:

            root_parts = root.split("/")
            if root_parts[-1] == "wins":
                continue

            if file.endswith('.csv'):
                dir_type = os.path.basename(root)
                config = configurations.get(dir_type)

                # Read the CSV file
                csv_path = os.path.join(root, file)
                df_table = pd.read_csv(csv_path, delimiter=',', index_col=0)

                # Plotting
                colors = ["red", "white", "green"]
                cmap = LinearSegmentedColormap.from_list("custom", colors, N=100)

                if file[:len("by-range")] == "by-range":
                    plt.figure(figsize=(14, 6))
                else:
                    plt.figure(figsize=(10, 6))

                formatted_data = df_table.applymap('{:.0f}'.format)

                ax = sns.heatmap(df_table, annot=True, cmap=cmap, center=config['cmap_center'], vmin=config['vmin'],
                                 vmax=config['vmax'], fmt='g')

                # Setting x and y axis labels
                ax.set_ylabel("Threshold")

                # Get the current tick labels
                xlabels = [label.get_text() for label in ax.get_xticklabels()]
                ylabels = [label.get_text() for label in ax.get_yticklabels()]

                # Format the labels
                ylabels = ["{:.1f}%".format(float(label.replace("thld:", "")) * 100) for label in ylabels]
                if xlabels[0][:len("limit")] == "limit":
                    ax.set_xlabel("Limit")
                    xlabels = [label.replace("limit:", "") for label in xlabels]
                else:
                    ax.set_xlabel("Range")
                    xlabels = [label.replace("rng:", "") for label in xlabels]

                xlabels = [str(int(float(label))) for label in xlabels]

                ax.set_xticklabels(xlabels)
                ax.set_yticklabels(ylabels)

                plt.title(file.replace("_", " ").replace(".csv", "").title())
                plt.tight_layout()

                # Construct the output path
                relative_path = os.path.relpath(csv_path, base_dir)
                svg_path = os.path.join(root_dir, './svg', os.path.splitext(relative_path)[0] + ".svg")

                # Create directories if they don't exist
                os.makedirs(os.path.dirname(svg_path), exist_ok=True)

                # Save the SVG file
                plt.savefig(svg_path, format="svg")
                plt.close()
                print(svg_path)


if __name__ == '__main__':
    generate_heatmaps("/Users/ruae/Projects/strategy-evaluator/output_7")
    generate_heatmaps("/Users/ruae/Projects/strategy-evaluator/output_14")
    generate_heatmaps("/Users/ruae/Projects/strategy-evaluator/output_21")
