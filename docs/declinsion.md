As your Product Manager, I have converted the **PRD: Greek Case Master (MVP)** into a Markdown format suitable for a Google Doc.

You can copy and paste the content below directly into a new Google Doc.

---

# PRD: Greek Case Master (MVP)

## 1\. Executive Summary

**Objective:** Transition a Modern Greek learner from rote memorization of noun tables to instinctive application of cases (Nominative, Genitive, Accusative) based on sentence context. **Product Persona:** Specialized Modern Greek Grammar Analyst.

## 2\. User Profile & Problem Statement

**User:** English speakers learning Modern Greek who struggle with real-time declension.

**Problem:** Noun cases feel abstract and difficult to apply "on the fly" during sentence construction.

## 3\. Core Functional Requirements

### FR1: The "Syntax Prompt" Engine

The system presents a sentence in English with a missing Greek noun component.

* **Input:** User provides the **Greek Article \+ Declined Noun**.

**Validation:** Match against the target declined form (e.g., matching "τον δάσκαλο" for an accusative masculine object).

### FR2: Automated Grammar Breakdown

Upon every answer, the system must provide a breakdown following the "Modern Greek Grammar Analyst" instructions:

**Translation:** Accurate translation of the full sentence.

**Syntactic Role:** Identification of the noun's function (e.g., Subject \= Nominative, Direct Object \= Accusative).

**Morphology:** Explanation of the specific gender and declension pattern (e.g., how "-ος" becomes "-ου" in the genitive).

### FR3: Level-Based Difficulty Scaling

**Phase 1 (Direct Objects):** Focus on the **Accusative** case after transitive verbs like *βλέπω* (I see) or *ψάχνω* (I look for).

**Phase 2 (Possession):** Focus on the **Genitive** case to show ownership (e.g., "The book of...").

**Phase 3 (Prepositions):** Practice cases required after common prepositions such as *σε* (to/in), *από* (from), or *για* (for).

## 4\. Technical Constraints & Design Rules

**The Article-Noun Unit:** Users must practice the article and noun together, as the article (e.g., *ο, του, τον*) is a primary indicator of case and gender.

**Brevity:** Feedback must be strictly analytical and avoid conversational filler or enthusiastic exclamations.

## 5\. PM Clarification Questions for Next Sprint

1. **Input Method:** Full manual typing (best for spelling/accents) or Word Bank/Multiple Choice (best for speed)?  
2. **Number Scope:** Should the MVP start strictly with **Singular** nouns or mix in **Plurals** immediately?  
3. **Strictness:** Should a correct noun ending be marked "Wrong" if the accompanying article is incorrect?

---

### Reference Material for Development:

**Nouns & Articles:** Detailed forms can be found in reference guides such as *Greek: An Essential Grammar of the Modern Language*.

* **Case Functions:**

**Nominative:** Subject of the verb.

**Accusative:** Direct object or object of most prepositions.

**Genitive:** Possession or specific prepositional use.  
