import pandas as pd
import matplotlib.pyplot as plt


df = pd.read_csv("loss.csv")


## 📊 Plot 1: Actor and Critic Loss Over Steps
plt.figure(figsize=(12, 6))
plt.plot(df.index, df["actor_loss"], label="Actor Loss", color="#1f77b4", marker="o", markersize=3, alpha=0.7)
plt.plot(df.index, df["critic_loss"], label="Critic Loss", color="#ff7f0e", marker="s", markersize=3, alpha=0.7)

plt.title("Actor vs Critic Loss Over Update Steps", fontsize=14)
plt.xlabel("Update Step")
plt.ylabel("Loss Value")
plt.grid(True, linestyle="--", alpha=0.6)
plt.legend()
plt.show()


## 📉 Plot 2: Smoothed Total Loss

plt.figure(figsize=(12, 6))
# Calculate and plot the rolling average of the Total Loss
df["total_loss"].rolling(window=3).mean().plot(
    label="Total Loss (smoothed, window=3)", color="#2ca02c", linestyle='-', linewidth=2
)

plt.title("Smoothed Total Loss (Actor + Critic)", fontsize=14)
plt.xlabel("Update Step")
plt.ylabel("Total Loss (Rolling Mean, window=3)")
plt.grid(True, linestyle="--", alpha=0.6)
plt.legend()
plt.show()


## 🔗 Plot 3: Actor Loss vs Critic Loss Relationship


plt.figure(figsize=(8, 6))
plt.scatter(df["critic_loss"], df["actor_loss"], color="#d62728", alpha=0.7, edgecolor="k")

plt.title("Critic Loss vs Actor Loss", fontsize=14)
plt.xlabel("Critic Loss")
plt.ylabel("Actor Loss")
plt.grid(True, linestyle="--", alpha=0.6)
plt.show()