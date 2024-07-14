import React, { ChangeEvent } from 'react';
import { FieldSet, InlineField, InlineFieldRow, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions } from './types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const {
    onOptionsChange,
    options,
    options: { jsonData },
  } = props;

  const onPlaystationIPChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      playstationIP: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  const { playstationIP } = jsonData;

  return (
    <FieldSet label="Connection">
      <InlineFieldRow>
        <InlineField label="Playstation IP" labelWidth={20} tooltip="IPv4 only for now">
          <Input
            width={20}
            data-testid="playstationIP"
            required
            value={playstationIP}
            autoComplete="off"
            placeholder="192.168.1.x"
            onChange={onPlaystationIPChange}
            css={undefined}
          />
        </InlineField>
      </InlineFieldRow>
    </FieldSet>
  );
}
