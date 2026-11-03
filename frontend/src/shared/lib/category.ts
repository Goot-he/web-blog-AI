export const OTHER_CATEGORY_NAME = '其它'

export const getCategoryName = (categoryID?: number, categoryNames?: Map<number, string>) => {
  if (!categoryID) return OTHER_CATEGORY_NAME
  return categoryNames?.get(categoryID) ?? OTHER_CATEGORY_NAME
}
