export const CAMPUS_SLUGS = [
  "internships",
  "postgraduate",
  "second-hand",
  "courses",
  "daily",
];

export function campusCategories(site) {
  return CAMPUS_SLUGS.map((slug) =>
    site.categories.find(
      (category) => category.slug === slug && !category.parent_category_id
    )
  ).filter(Boolean);
}
