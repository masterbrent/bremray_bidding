import { test, expect } from '@playwright/test';

test('phases should display and update', async ({ page }) => {
  // Navigate to templates page
  await page.goto('http://localhost:5174/#/templates');
  
  // Wait for templates to load
  await page.waitForSelector('.template-card');
  
  // Click on Sunroom template
  await page.click('text=Sunroom');
  
  // Wait for template details
  await page.waitForSelector('.template-details');
  
  // Check if phases section exists
  const phasesSection = await page.locator('.section:has-text("Phases")');
  
  // Take a screenshot
  await page.screenshot({ path: 'phases-test.png', fullPage: true });
  
  console.log('Test completed - check phases-test.png');
});
