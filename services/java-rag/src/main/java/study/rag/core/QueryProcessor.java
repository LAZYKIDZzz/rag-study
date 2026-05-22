package study.rag.core;

import org.springframework.stereotype.Component;
import study.rag.core.model.ChatSession;
import study.rag.core.model.UserMemory;

import java.util.ArrayList;
import java.util.List;

@Component
public class QueryProcessor {
  public String rewrite(String question, ChatSession session, UserMemory memory) {
    List<String> parts = new ArrayList<>();
    if (!memory.facts().isEmpty()) {
      List<String> facts = memory.facts();
      parts.add("User facts: " + String.join("; ", facts.subList(Math.max(0, facts.size() - 3), facts.size())));
    }
    if (session != null) {
      session.messages().stream()
          .filter(message -> message.role().equals("user"))
          .reduce((first, second) -> second)
          .ifPresent(message -> parts.add("Previous question: " + message.content()));
    }
    parts.add("Current question: " + question);
    return String.join("\n", parts);
  }
}
