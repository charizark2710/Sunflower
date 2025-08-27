import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import numpy as np

# Load data
df = pd.read_csv("output.csv")
df["ic"] = np.log10(df["ic"])

df["x_formula"] = (
    df["bundleSize"] + df["totalCallCount"] *
    (
        df["averageCyclomaticComplexity"]
        + (df["totalConditionalComplexity"] * 0.5)
        + df["totalLoopComplexity"]
        + (df["totalRecursionDepth"] * 1.75)
    )
)

# Compute correlation matrix
corr = df.corr(numeric_only=True)

# Create figure with 2 subplots (vertical layout)
fig, axes = plt.subplots(
    2, 1, figsize=(18, 9),
    gridspec_kw={'height_ratios': [2, 1]}  # heatmap taller
)

# --- Heatmap ---
corr = df.corr(numeric_only=True)
sns.heatmap(corr, annot=True, cmap="coolwarm", fmt=".2f", cbar=True, ax=axes[0])
axes[0].set_title("Feature Correlation Heatmap", fontsize=14)

# --- Line Graph for IC ---
sampled = df.iloc[::50]
axes[1].plot(sampled.index, sampled["ic"], color="blue", linestyle="-", marker="o",
             markersize=4, markerfacecolor="black", label="IC")
axes[1].set_title("Instruction Count (IC) per Sample", fontsize=14)
axes[1].set_xlabel("Sample Index")
axes[1].set_ylabel("IC")
axes[1].grid(True, linestyle="--", alpha=0.6)
axes[1].legend()

plt.subplots_adjust(hspace=0.4)  # increase spacing
plt.show()