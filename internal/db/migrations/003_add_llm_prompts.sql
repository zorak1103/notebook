-- Add LLM prompt configuration keys with default values

INSERT INTO config (key, value) VALUES ('llm_prompt_summary',
'Create a concise meeting summary (3-5 sentences) based on the notes below. Capture key discussion points, decisions, and action items.

IMPORTANT:
- Write the summary in the same language as the notes. Do not translate.
- Provide ONLY the final summary text. Do NOT offer multiple options or ask the user to choose.
- Output the summary directly without any introduction, preamble, or explanation.

Meeting: {{subject}}
Date: {{date}}
Participants: {{participants}}

Notes:
{{notes}}');

INSERT INTO config (key, value) VALUES ('llm_prompt_enhance',
'Improve the following meeting note: fix grammar and spelling, improve clarity and structure. Preserve all original information and meaning.

IMPORTANT:
- Write in the same language as the original note. Do not translate.
- Provide ONLY the improved note text. Do NOT offer multiple versions or ask the user to choose.
- Output the enhanced note directly without any introduction, preamble, or explanation.

Note:
{{content}}');
