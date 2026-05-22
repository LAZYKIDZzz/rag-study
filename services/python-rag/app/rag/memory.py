from app.rag.models import UserMemory


class MemoryStore:
    def __init__(self) -> None:
        self._memories: dict[str, UserMemory] = {}

    def get(self, user_id: str) -> UserMemory:
        if user_id not in self._memories:
            self._memories[user_id] = UserMemory(user_id=user_id)
        return self._memories[user_id]

    def add_fact(self, user_id: str, fact: str) -> UserMemory:
        memory = self.get(user_id)
        cleaned = fact.strip()
        if cleaned and cleaned not in memory.facts:
            memory.facts.append(cleaned)
        return memory

    def update_summary(self, user_id: str, latest_user_message: str) -> UserMemory:
        memory = self.get(user_id)
        if memory.recent_summary:
            memory.recent_summary = f"{memory.recent_summary} | Last asked: {latest_user_message[:120]}"
        else:
            memory.recent_summary = f"Last asked: {latest_user_message[:120]}"
        return memory
