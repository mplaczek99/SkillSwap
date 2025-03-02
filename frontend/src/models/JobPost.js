// JobPost.js - Model for job postings
export default class JobPost {
  constructor(data = {}) {
    this.id = data.id || null;
    this.title = data.title || "";
    this.company = data.company || "";
    this.location = data.location || "";
    this.description = data.description || "";
    this.skillsRequired = data.skillsRequired || [];
    this.experienceLevel = data.experienceLevel || "Entry"; // Entry, Mid, Senior
    this.jobType = data.jobType || "Full-time"; // Full-time, Part-time, Contract
    this.salaryRange = data.salaryRange || "";
    this.contactEmail = data.contactEmail || "";
    this.postedByUserID = data.postedByUserID || null;
    this.postedByName = data.postedByName || "";
    this.createdAt = data.createdAt ? new Date(data.createdAt) : new Date();
    this.updatedAt = data.updatedAt ? new Date(data.updatedAt) : new Date();
  }

  // Helper method to format the creation date
  formattedDate() {
    return this.createdAt.toLocaleDateString(undefined, {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  }

  // Helper method to get days since posting
  daysSincePosting() {
    const now = new Date();
    const diffTime = Math.abs(now - this.createdAt);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  }

  // Helper method to get an array of skills
  skillsArray() {
    if (Array.isArray(this.skillsRequired)) {
      return this.skillsRequired;
    } else if (typeof this.skillsRequired === "string") {
      return this.skillsRequired.split(",").map((skill) => skill.trim());
    }
    return [];
  }

  // Convert to a simple object
  toJSON() {
    return {
      id: this.id,
      title: this.title,
      company: this.company,
      location: this.location,
      description: this.description,
      skillsRequired: this.skillsRequired,
      experienceLevel: this.experienceLevel,
      jobType: this.jobType,
      salaryRange: this.salaryRange,
      contactEmail: this.contactEmail,
      postedByUserID: this.postedByUserID,
      postedByName: this.postedByName,
      createdAt: this.createdAt,
      updatedAt: this.updatedAt,
    };
  }
}
