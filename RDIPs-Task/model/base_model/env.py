from transformers import RobertaTokenizer, RobertaModel
from utils.common import DEVICE

from gymnasium.envs.registration import register
import gymnasium as gym
from gymnasium import spaces
import numpy as np
import psutil
import time

CPU_RATE = 3.5e9

register(
    id='task-eval-v0',
    entry_point='v0_task_eval_env:TaskEvalEnv',
)

SAVE_DIR = "./saved_model"
LR = 1e-4

class TaskEvalEnv(gym.Env):
    def __init__(self):
        super(TaskEvalEnv, self).__init__()
        psutil.cpu_percent(interval=5)
        self.memory_total = psutil.virtual_memory().total
        self.action_space = spaces.Discrete(2)  # Skip or Execute
        self.observation_space = spaces.Box(
            low=np.array([0.0, 0.0, 0.0]),
            high=np.array([1.0, 100.0, 100.0]),
            shape=(3,),
            dtype=np.float64
        )

        # placeholders for context updated externally
        self.result = {}
        self.task = {}

    def update_context(self, result=None, task=None):
        """
        Called externally before env.step(action).
        Stores the latest result and task context.
        """
        if task is not None:
            self.task = task
        if result is not None:
            self.result = result

    def step(self, action, mem_usage, cpu_usage):
        certainty = self.task["certainty"]
        result = self.result
        reward = 0.0
        executed = int(action) == 1

        if mem_usage is None or cpu_usage is None:
            mem_usage = psutil.virtual_memory().used
            cpu_usage = psutil.cpu_percent(interval=0)
            mem_percent = mem_usage / self.memory_total
            cpu_percent = cpu_usage / 100.0
        else:
            mem_percent = abs(mem_usage) / self.memory_total
            cpu_percent = abs(cpu_usage) / 100.0
        
        if executed:
            percent_deduction = 1.0
            if mem_percent > 0.9:
                percent_deduction *= 2.0
            elif mem_percent > 0.7:
                percent_deduction *= 1.5
            elif mem_percent > 0.5:
                percent_deduction *= 1.2

            if cpu_percent > 0.9:
                percent_deduction *= 2.0
            elif cpu_percent > 0.7:
                percent_deduction *= 1.5
            elif cpu_percent > 0.5:
                percent_deduction *= 1.2
        else:
            percent_deduction = -1.0
            if mem_percent < 0.1:
                percent_deduction *= 2.5
            elif mem_percent < 0.3:
                percent_deduction *= 2.0
            elif mem_percent < 0.5:
                percent_deduction *= 1.2

            if cpu_percent < 0.1:
                percent_deduction *= 2.5
            elif cpu_percent < 0.3:
                percent_deduction *= 2.0
            elif cpu_percent < 0.5:
                percent_deduction *= 1.2

            if certainty < 0.7 and percent_deduction > 0:
                percent_deduction *= 2.0

        # Reward logic
        if executed:
            if result is None:
                reward = -5.0
            elif certainty > 0.7:
                reward = 2.0
            elif 0.5 <= certainty <= 0.7:
                reward = 1.0
            else:  # certainty < 0.5
                reward = 0.5
        else:
            if certainty > 0.7:
                reward = 2.0
            elif 0.5 <= certainty <= 0.7:
                reward = 1.0
            else:  # certainty < 0.5
                reward = 0.5
        done = True
        obs = [mem_percent, cpu_percent]
        return obs, reward / percent_deduction, done
    def _get_state(self):
        """Return [mem_usage_ratio, cpu_usage%]."""
        cpu = psutil.cpu_percent(interval=0) 
        cpu_usage = cpu / 100
        mem_usage = psutil.virtual_memory().used / self.memory_total
        return np.array([mem_usage, cpu_usage], dtype=np.float32)

    def reset(self, *, seed=None, options=None):
        super().reset(seed=seed)
        self.result = {}
        self.task = {}
        time.sleep(1)
        return self._get_state()
