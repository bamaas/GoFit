import { writable, derived } from 'svelte/store';

export type Metadata = {
  total_records: number,
  current_page: number,
  page_size: number,
  first_page: number,
  last_page: number
}

export type CheckIn = {
  uuid: string;
  datetime: number;
  weight: number;
  moving_average: number;
  weight_difference: number;
};

export type ApiResponse = {
  data: CheckIn[];
  metadata: Metadata;
}

// Store for your data.
export const apiData = writable({} as ApiResponse);

// Data transformation
export const checkIns = derived(apiData, ($apiData) => {
    if ($apiData.data == undefined) {
        return [];
    }
    return $apiData.data;
});