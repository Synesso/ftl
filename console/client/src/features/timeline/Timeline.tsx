import { Timestamp } from '@bufbuild/protobuf'
import React, { useContext, useEffect, useState } from 'react'
import { useClient } from '../../hooks/use-client.ts'
import { ConsoleService } from '../../protos/xyz/block/ftl/v1/console/console_connect.ts'
import { Module, StreamTimelineResponse } from '../../protos/xyz/block/ftl/v1/console/console_pb.ts'
import { SidePanelContext } from '../../providers/side-panel-provider.tsx'
import { classNames } from '../../utils/react.utils.ts'
import { TimelineCall } from './TimelineCall.tsx'
import { TimelineDeployment } from './TimelineDeployment.tsx'
import { TimelineEventDetailCall } from './TimelineEventDetailCall.tsx'
import { TimelineEventDetailDeployment } from './TimelineEventDetailDeployment.tsx'
import { TimelineEventDetailLog } from './TimelineEventDetailLog.tsx'
import { TimelineLog } from './TimelineLog.tsx'

type Props = {
  module?: Module | null
}

export const Timeline: React.FC<Props> = ({ module }) => {
  const client = useClient(ConsoleService)
  const { openPanel, closePanel } = useContext(SidePanelContext)
  const [ entries, setEntries ] = useState<StreamTimelineResponse[]>([])
  const [ selectedEntry, setSelectedEntry ] = useState<StreamTimelineResponse | null>(null)

  useEffect(() => {
    const abortController = new AbortController()

    async function streamTimeline() {
      const afterTime = new Date()
      afterTime.setHours(afterTime.getHours() - 1)

      for await (const response of client.streamTimeline(
        { afterTime: Timestamp.fromDate(afterTime), deploymentName: module?.deploymentName },
        { signal: abortController.signal })
      ) {
        if (response.entry) {
          setEntries(prevEntries => [ response, ...prevEntries ])
        }
      }
    }

    streamTimeline()
    return () => {
      abortController.abort()
    }
  }, [ client, module ])

  const handleEntryClicked = entry => {
    if (selectedEntry === entry) {
      setSelectedEntry(null)
      closePanel()
      return
    }

    switch (entry.entry?.case) {
      case 'call':
        openPanel(<TimelineEventDetailCall call={entry.entry.value} />)
        break
      case 'log':
        openPanel(<TimelineEventDetailLog log={entry.entry.value} />)
        break
      case 'deployment':
        openPanel(<TimelineEventDetailDeployment deployment={entry.entry.value} />)
        break
      default: break
    }
    setSelectedEntry(entry)
  }

  return (
    <>
      <ul role='list' className='space-y-4'>
        {entries.map((entry, index) => (
          <li key={index} className='relative flex gap-x-4' onClick={() => handleEntryClicked(entry)}>
            <div className={classNames(index === entries.length - 1 ? 'h-6' : '-bottom-6', 'absolute left-0 top-0 flex w-6 justify-center')}>
              <div className='w-px bg-gray-200 dark:bg-gray-600' />
            </div>

            {(() => {
              switch (entry.entry?.case) {
                case 'call': return <TimelineCall call={entry.entry.value} />
                case 'log': return <TimelineLog log={entry.entry.value} />
                case 'deployment': return <TimelineDeployment
                  deployment={entry.entry.value}
                  timestamp={entry.timeStamp}
                />
                default: return <></>
              }
            })()}
          </li>
        ))}
      </ul>
    </>
  )
}