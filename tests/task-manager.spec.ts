import { test, expect } from '@playwright/test';

test.beforeEach(async ({ page, request }) => {
  await request.post('/api/reset');
  await page.goto('/');
  await expect(page.locator('#statTotal')).toHaveText('0');
});

test('page loads with correct title', async ({ page }) => {
  await expect(page).toHaveTitle('Task Manager');
  await expect(page.getByText('Get things done.')).toBeVisible();
});

test('add a task', async ({ page }) => {
  await page.getByPlaceholder('Add a new task...').fill('Buy groceries');
  await page.getByRole('button', { name: 'Add' }).click();

  await expect(page.getByText('Buy groceries')).toBeVisible();
  await expect(page.locator('#statTotal')).toHaveText('1');
  await expect(page.locator('#statPending')).toHaveText('1');
});

test('add task with Enter key', async ({ page }) => {
  await page.getByPlaceholder('Add a new task...').fill('Write tests');
  await page.keyboard.press('Enter');

  await expect(page.getByText('Write tests')).toBeVisible();
});

test('mark task as complete', async ({ page }) => {
  await page.getByPlaceholder('Add a new task...').fill('Fix bug');
  await page.getByRole('button', { name: 'Add' }).click();

  await page.locator('.check-btn').first().click();

  await expect(page.locator('.task-item.done')).toBeVisible();
  await expect(page.locator('#statDone')).toHaveText('1');
  await expect(page.locator('#statPending')).toHaveText('0');
});

test('toggle task back to incomplete', async ({ page }) => {
  await page.getByPlaceholder('Add a new task...').fill('Deploy app');
  await page.getByRole('button', { name: 'Add' }).click();

  await page.locator('.check-btn').first().click();
  await expect(page.locator('.task-item.done')).toBeVisible();

  await page.locator('.check-btn').first().click();
  await expect(page.locator('.task-item.done')).not.toBeVisible();
  await expect(page.locator('#statPending')).toHaveText('1');
});

test('delete a task', async ({ page }) => {
  await page.getByPlaceholder('Add a new task...').fill('Task to delete');
  await page.getByRole('button', { name: 'Add' }).click();
  await expect(page.getByText('Task to delete')).toBeVisible();

  await page.locator('.del-btn').first().click();

  await expect(page.getByText('Task to delete')).not.toBeVisible();
  await expect(page.locator('#statTotal')).toHaveText('0');
});

test('filter pending tasks', async ({ page }) => {
  await page.getByPlaceholder('Add a new task...').fill('Task one');
  await page.getByRole('button', { name: 'Add' }).click();
  await page.getByPlaceholder('Add a new task...').fill('Task two');
  await page.getByRole('button', { name: 'Add' }).click();
  await expect(page.getByText('Task two', { exact: true })).toBeVisible();

  await page.locator('.check-btn').first().click();

  await page.getByRole('button', { name: 'Pending' }).click();
  await expect(page.getByText('Task two', { exact: true })).toBeVisible();
  await expect(page.getByText('Task one', { exact: true })).not.toBeVisible();
});

test('filter done tasks', async ({ page }) => {
  await page.getByPlaceholder('Add a new task...').fill('Task A');
  await page.getByRole('button', { name: 'Add' }).click();
  await page.getByPlaceholder('Add a new task...').fill('Task B');
  await page.getByRole('button', { name: 'Add' }).click();
  await expect(page.getByText('Task B', { exact: true })).toBeVisible();

  await page.locator('.check-btn').first().click();

  await page.getByRole('button', { name: 'Done' }).click();
  await expect(page.getByText('Task A', { exact: true })).toBeVisible();
  await expect(page.getByText('Task B', { exact: true })).not.toBeVisible();
});

test('empty state shows when no tasks', async ({ page }) => {
  await expect(page.locator('#emptyMsg')).toBeVisible();
  await expect(page.locator('#emptyMsg')).toHaveText('No tasks yet. Add one above!');
});

test('stats update correctly', async ({ page }) => {
  await expect(page.locator('#statTotal')).toHaveText('0');

  await page.getByPlaceholder('Add a new task...').fill('Task 1');
  await page.getByRole('button', { name: 'Add' }).click();
  await expect(page.getByText('Task 1', { exact: true })).toBeVisible();
  await page.getByPlaceholder('Add a new task...').fill('Task 2');
  await page.getByRole('button', { name: 'Add' }).click();
  await expect(page.getByText('Task 2', { exact: true })).toBeVisible();

  await expect(page.locator('#statTotal')).toHaveText('2');
  await expect(page.locator('#statPending')).toHaveText('2');
  await expect(page.locator('#statDone')).toHaveText('0');

  await page.locator('.check-btn').first().click();

  await expect(page.locator('#statDone')).toHaveText('1');
  await expect(page.locator('#statPending')).toHaveText('1');
});
