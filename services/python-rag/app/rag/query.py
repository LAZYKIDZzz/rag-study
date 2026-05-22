from app.rag.models import ChatSession, UserMemory


class QueryProcessor:
    def rewrite(self, question: str, session: ChatSession | None, memory: UserMemory) -> str:
        parts: list[str] = []
        if memory.facts:
            parts.append("User facts: " + "; ".join(memory.facts[-3:]))

        previous_user_messages = [
            message.content for message in (session.messages if session else []) if message.role == "user"
        ]
        if previous_user_messages:
            parts.append("Previous question: " + previous_user_messages[-1])

        parts.append("Current question: " + question)
        return "\n".join(parts)
