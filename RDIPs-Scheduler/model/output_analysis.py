import pandas as pd
import matplotlib.pyplot as plt

# --- Load dataset ---
df = pd.read_csv("history.csv")

# --- Line chart: Train vs Eval loss ---
plt.figure(figsize=(12, 6))
plt.plot(df.index, df["train_loss"], label="Train Loss", color="blue", marker="o", markersize=3, alpha=0.7)
plt.plot(df.index, df["eval_loss"], label="Eval Loss", color="red", marker="s", markersize=3, alpha=0.7)

plt.title("Training vs Evaluation Loss over Steps", fontsize=14)
plt.xlabel("Step")
plt.ylabel("Loss")
plt.grid(True, linestyle="--", alpha=0.6)
plt.legend()
plt.show()

# --- Smoothed curves using rolling average ---
plt.figure(figsize=(12, 6))
df["train_loss"].rolling(window=5).mean().plot(label="Train Loss (smoothed)", color="blue")
df["eval_loss"].rolling(window=5).mean().plot(label="Eval Loss (smoothed)", color="red")

plt.title("Smoothed Training vs Evaluation Loss", fontsize=14)
plt.xlabel("Step")
plt.ylabel("Loss (Rolling Mean, window=5)")
plt.grid(True, linestyle="--", alpha=0.6)
plt.legend()
plt.show()

# --- Scatter plot: Train vs Eval loss relationship ---
plt.figure(figsize=(8, 6))
plt.scatter(df["train_loss"], df["eval_loss"], color="purple", alpha=0.7, edgecolor="k")

plt.title("Train Loss vs Eval Loss", fontsize=14)
plt.xlabel("Train Loss")
plt.ylabel("Eval Loss")
plt.grid(True, linestyle="--", alpha=0.6)
plt.show()
