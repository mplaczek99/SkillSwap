import nlp from "compromise";

const iconCache = new Map();

/**
 * fetchDynamicIcon returns a Font Awesome icon name for the given skill.
 * It queries Iconify’s API (using the "fa-solid" prefix) and returns a fallback if needed.
 */
export async function fetchDynamicIcon(skillName) {
  const normalizedName = skillName.trim().toLowerCase();
  if (iconCache.has(normalizedName)) {
    return iconCache.get(normalizedName);
  }

  // Extract key nouns from the skill name.
  const doc = nlp(skillName);
  const topics = doc.nouns().out("array");
  const query = topics.length ? topics[0] : skillName;

  let iconName = "cog"; // default fallback
  try {
    const response = await fetch(
      `https://api.iconify.design/search?query=${encodeURIComponent(query)}&prefix=fa-solid`,
    );
    if (response.ok) {
      const data = await response.json();
      if (data && data.results && data.results.length > 0) {
        iconName = data.results[0].icon;
        // Remove any "fa-" prefix for consistency.
        if (iconName.startsWith("fa-")) {
          iconName = iconName.substring(3);
        }
      }
    } else {
      console.error("Icon API response not OK:", response.status);
    }
  } catch (error) {
    console.error("Dynamic icon lookup failed:", error);
  }

  iconCache.set(normalizedName, iconName);
  return iconName;
}
