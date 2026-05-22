from app.rag.models import SearchResult, UserMemory


class ExtractiveAnswerGenerator:
    """Builds a local answer from retrieved chunks so the pipeline works offline."""

    def generate(self, question: str, results: list[SearchResult], memory: UserMemory) -> str:
        if not results:
            return (
                "I could not find matching knowledge base content yet. "
                "Add a document, index it, then ask again."
            )

        facts = ""
        if memory.facts:
            facts = " I also considered your saved facts: " + "; ".join(memory.facts[-3:]) + "."

        lead = f"Based on the retrieved knowledge, here is the most relevant answer to: {question}"
        snippets = []
        for index, result in enumerate(results[:3], start=1):
            content = result.chunk.content.strip()
            if len(content) > 420:
                content = content[:417].rstrip() + "..."
            snippets.append(f"[{index}] {content}")
        return lead + facts + "\n\n" + "\n\n".join(snippets)
