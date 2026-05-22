package study.rag.core.model;

import java.util.ArrayList;
import java.util.List;

public class UserMemory {
  private final String userId;
  private final List<String> facts = new ArrayList<>();
  private String recentSummary = "";

  public UserMemory(String userId) {
    this.userId = userId;
  }

  public String userId() {
    return userId;
  }

  public List<String> facts() {
    return List.copyOf(facts);
  }

  public String recentSummary() {
    return recentSummary;
  }

  public void addFact(String fact) {
    if (!fact.isBlank() && !facts.contains(fact)) {
      facts.add(fact);
    }
  }

  public void updateSummary(String latestQuestion) {
    String addition = "Last asked: " + latestQuestion.substring(0, Math.min(120, latestQuestion.length()));
    recentSummary = recentSummary.isBlank() ? addition : recentSummary + " | " + addition;
  }
}
